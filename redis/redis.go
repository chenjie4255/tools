package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chenjie4255/tools/errcode"
	"github.com/chenjie4255/tools/gor"
	"github.com/chenjie4255/errors"
	"github.com/gomodule/redigo/redis"
)

// NewPool 生成redis.Pool
func NewPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewDB(host, password string, dbNum int) DB {
	pool := NewPool(host, password, dbNum)
	return &db{pool}
}

func NewDBFromPool(pool *redis.Pool) DB {
	return &db{pool}
}

func FlushDB(host, password string, dbNum int) error {
	pool := NewPool(host, password, dbNum)
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHDB")
	return err
}

type db struct {
	pool *redis.Pool
}

func (d *db) Get(key string, v interface{}) error {
	conn := d.pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return errors.NewWithTag("key doesn't exists", errcode.ResNotFound)
		}
		return err
	}

	return json.Unmarshal(data, &v)
}

func (d *db) Set(key string, value interface{}, expires int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return d.set(key, data, expires)
}

func (d *db) Exists(key string) (bool, error) {
	conn := d.pool.Get()
	defer conn.Close()
	result, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

const setNotExistsScript = `if redis.call("SETNX", KEYS[1], ARGV[1]) == 1 then
return redis.call("EXPIRE", KEYS[1], ARGV[2])
else
return 0
end`

func (d *db) setNotExists(key string, value interface{}, expires int) error {
	if expires <= 0 {
		panic("expires shoudl be greater than zero")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	scr := redis.NewScript(1, setNotExistsScript)
	conn := d.pool.Get()
	defer conn.Close()

	ret, err := redis.Int(scr.Do(conn, key, data, expires))
	if err != nil {
		return err
	}

	if ret != 1 {
		return errors.NewWithTag("resources exists", errcode.ResExisted)
	}

	return nil
}

func (d *db) SetNotExists(key string, value interface{}, expires int) error {
	if expires > 0 {
		return d.setNotExists(key, value, expires)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	conn := d.pool.Get()
	defer conn.Close()
	result, err := redis.Int(conn.Do("SETNX", key, data))
	if err != nil {
		return err
	}

	if result != 1 {
		return errors.NewWithTag("resources exists", errcode.ResExisted)
	}

	return nil
}

func (d *db) set(key string, data []byte, expires int) error {
	conn := d.pool.Get()
	defer conn.Close()
	if expires > 0 {
		_, err := conn.Do("SET", key, data, "EX", expires)
		return err
	}

	_, err := conn.Do("SET", key, data)
	return err

}

func (d *db) Del(key string) error {
	conn := d.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

func (d *db) DelByKeys(keys string) error {
	conn := d.pool.Get()
	defer conn.Close()

	items, err := redis.Strings(conn.Do("KEYS", keys))
	if err != nil {
		return err
	}

	ifaces := []interface{}{}
	for _, item := range items {
		ifaces = append(ifaces, item)
	}

	_, err = conn.Do("DEL", ifaces...)
	return err
}

func (d *db) DelKeysByScan(keyPattern string) error {
	conn := d.pool.Get()
	defer conn.Close()

	iter := 0
	var keys []string
	for {
		if arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", keyPattern)); err != nil {
			return err
		} else {
			iter, _ = redis.Int(arr[0], nil)
			keys, _ = redis.Strings(arr[1], nil)
		}

		delKeys := []interface{}{}
		for i := range keys {
			delKeys = append(delKeys, keys[i])
		}

		if len(delKeys) > 0 {
			if _, err := conn.Do("DEL", delKeys...); err != nil {
				return err
			}
		}

		if iter == 0 {
			break
		}
	}

	return nil
}

func (d *db) DelKeys(keys []string) error {
	conn := d.pool.Get()
	defer conn.Close()

	ifaces := []interface{}{}
	for _, item := range keys {
		ifaces = append(ifaces, item)
	}

	_, err := conn.Do("DEL", ifaces...)
	return err
}

const delKeyValueScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("del",KEYS[1])
else
return 0
end`

func (d *db) DelKeyForValue(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return d.DelKeyForBytes(key, jsonData)
}

func (d *db) DelKeyForBytes(key string, value []byte) error {
	scr := redis.NewScript(1, delKeyValueScript)
	conn := d.pool.Get()
	defer conn.Close()

	ret, err := redis.Int(scr.Do(conn, key, value))
	if err != nil {
		return err
	}

	if ret == 0 {
		return errors.New("key does not exist or unmatched")
	}

	return nil
}

const replaceValueScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("set",KEYS[1], ARGV[2])
else
return "NOOK"
end`

func (d *db) ReplaceValue(key string, oldValue, newValue interface{}) error {
	oldData, err := json.Marshal(oldValue)
	if err != nil {
		return err
	}

	newData, err := json.Marshal(newValue)
	if err != nil {
		return err
	}

	scr := redis.NewScript(1, replaceValueScript)
	conn := d.pool.Get()
	defer conn.Close()

	ret, err := redis.String(scr.Do(conn, key, oldData, newData))
	if err != nil {
		return err
	}

	if ret != "OK" {
		return errors.New("key does not exist or old value is not matched")
	}

	return nil
}

const cmpSetScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("set",KEYS[2], ARGV[2])
else
return "NOOK"
end`

func (d *db) CmpAndSet(cmpKey string, cmpValue interface{}, setKey string, setValue interface{}) error {
	cmpData, err := json.Marshal(cmpValue)
	if err != nil {
		return err
	}

	setData, err := json.Marshal(setValue)
	if err != nil {
		return err
	}

	scr := redis.NewScript(2, cmpSetScript)
	conn := d.pool.Get()
	defer conn.Close()

	ret, err := redis.String(scr.Do(conn, cmpKey, setKey, cmpData, setData))
	if err != nil {
		return err
	}

	if ret != "OK" {
		return errors.New("key does not exist or old value is not matched")
	}

	return nil
}

const cmpGTDecrScript = `if tonumber(redis.call("GET",KEYS[1])) > tonumber(ARGV[1]) then
return redis.call("decr",KEYS[1])
else
return redis.error_reply("not exists or unmatch")
end`

func (d *db) CmpGTDecr(cmpKey string, greatThan int64) (int64, error) {
	scr := redis.NewScript(1, cmpGTDecrScript)
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Int64(scr.Do(conn, cmpKey, greatThan))
}

func (d *db) CacheGet(key string, v interface{}, fn FetchFunc, expires int) error {
	if err := d.Get(key, v); err == nil {
		// cache hit
		return nil
	}

	data, err := fn()
	if err != nil {
		return err
	}

	rawData, _ := json.Marshal(data)
	gor.RunWithRecover(func() {
		d.set(key, rawData, expires)
	})

	return json.Unmarshal(rawData, v)
}

func (d *db) CacheGetB(key string, v interface{}, fn FetchFunc, expires int) (bool, error) {
	if err := d.Get(key, v); err == nil {
		// cache hit
		return true, nil
	}

	data, err := fn()
	if err != nil {
		return false, err
	}

	rawData, _ := json.Marshal(data)
	gor.RunWithRecover(func() {
		d.set(key, rawData, expires)
	})

	return false, json.Unmarshal(rawData, v)
}

func (d *db) AddSortSetStr(key string, value string, sortKey int64) error {
	conn := d.pool.Get()
	defer conn.Close()

	result, err := redis.Int(conn.Do("ZADD", key, sortKey, value))
	if err != nil {
		return err
	}

	if result == 0 {
		return errors.NewWithTag("item exists", errcode.ResExisted)
	}

	return nil
}

func (d *db) GetSortSetCount(key string, sortKeyFrom, sortKeyTo int64) (int, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("ZCOUNT", key, sortKeyFrom, sortKeyTo))
}

func (d *db) GetSortSetRangeStr(key string, sortKeyFrom, sortKeyTo int64) ([]string, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("ZRANGEBYSCORE", key, sortKeyFrom, sortKeyTo))
}

func (d *db) RemoveSortSet(key string, sortKeyFrom, sortKeyTo int64) (int, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("ZREMRANGEBYSCORE", key, sortKeyFrom, sortKeyTo))
}

func (d *db) PushStringList(key string, value string, expires int) error {
	conn := d.pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		return err
	}

	result, err := redis.Int(conn.Do("EXPIRE", key, expires))
	if err != nil {
		return err
	}

	if result != 1 {
		return errors.NewWithTag("resources not exists", errcode.ResNotFound)
	}

	return nil
}

func (d *db) GetStringList(key string) ([]string, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("LRANGE", key, 0, -1))
}

func (d *db) TTL(key string) (int, error) {
	conn := d.pool.Get()
	defer conn.Close()

	result, err := redis.Int(conn.Do("TTL", key))
	if err != nil {
		return 0, err
	}

	if result == -2 {
		return 0, errors.NewWithTag("key does not exists", errcode.ResNotFound)
	}

	return result, nil
}

func (d *db) IncrByUint64(key string, step uint64) (uint64, error) {
	conn := d.pool.Get()
	defer conn.Close()

	result, err := redis.Uint64(conn.Do("INCRBY", key, step))
	if err != nil {
		return 0, err
	}

	return result, err
}

func (d *db) DecrByUint64(key string, step uint64) (uint64, error) {
	conn := d.pool.Get()
	defer conn.Close()

	result, err := redis.Uint64(conn.Do("DECRBY", key, step))
	if err != nil {
		return 0, err
	}

	return result, err
}

func (d *db) GetSetUInt64(key string, val uint64) (uint64, error) {
	conn := d.pool.Get()
	defer conn.Close()

	result, err := redis.Uint64(conn.Do("GETSET", key, val))
	if err != nil {
		return 0, err
	}

	return result, err
}

const incrToScript = `if redis.call("get",KEYS[1]) < ARGV[1] then
return redis.call("SET",KEYS[1], ARGV[1])
else
return "NOOK"
end`

func (d *db) IncrToUint64(key string, val uint64) (uint64, error) {
	scr := redis.NewScript(1, incrToScript)
	conn := d.pool.Get()
	defer conn.Close()

	ret, err := redis.String(scr.Do(conn, key, val))
	if err != nil {
		return 0, err
	}

	if ret != "OK" {
		return 0, nil
	}

	return val, nil
}

func (d *db) Time() (int64, error) {
	conn := d.pool.Get()
	defer conn.Close()

	ret, err := redis.Strings(conn.Do("TIME"))
	if err != nil {
		return 0, err
	}

	ts, err := strconv.ParseInt(ret[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return ts, nil
}

func MarshalJSONSlice(ifs []interface{}) ([][]byte, error) {
	ret := [][]byte{}
	for i := range ifs {
		data, err := json.Marshal(ifs[i])
		if err != nil {
			return nil, err
		}
		ret = append(ret, data)
	}

	return ret, nil
}

func UnMarshalJSONSlice(data [][]byte, output []interface{}) error {
	if len(data) != len(output) {
		return fmt.Errorf("the lenght of data and output is not equal %d != %d", len(data), len(output))
	}

	for i := range data {
		err := json.Unmarshal(data[i], output[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *db) SAdd(key string, values ...[]byte) error {
	conn := d.pool.Get()
	defer conn.Close()

	params := []interface{}{key}
	for i := range values {
		params = append(params, values[i])
	}

	_, err := redis.Int(conn.Do("SADD", params...))
	return err
}

func (d *db) SRandMember(key string, count int) ([][]byte, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("SRANDMEMBER", key, count))
}

func (d *db) SCard(key string) (int, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("SCARD", key))
}

func (d *db) Subscribe(channels []string, done chan bool, recv func(name string, data []byte)) error {
	conn := d.pool.Get()
	defer conn.Close()
	pcs := redis.PubSubConn{
		Conn: conn,
	}

	chs := []interface{}{}
	for i := range channels {
		chs = append(chs, channels[i])
	}

	if err := pcs.Subscribe(chs...); err != nil {
		return err
	}

	subCount := 0

	for {
		switch v := pcs.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message %s:\n", v.Channel)
			recv(v.Channel, v.Data)
		case error:
			fmt.Printf("error %s:\n", v)
			return v
		case redis.Subscription:
			fmt.Printf("Subscription %s: %s %d\n", v.Channel, v.Kind, v.Count)
			subCount++
			if subCount == len(channels) {
				done <- true
			}
		}
	}
}

func (d *db) Publish(channel string, data []byte) (int, error) {
	conn := d.pool.Get()
	defer conn.Close()

	count, err := redis.Int(conn.Do("PUBLISH", channel, data))
	return count, err
}
