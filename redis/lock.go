package redis

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/chenjie4255/errors"
	"github.com/gomodule/redigo/redis"
)

func randomValue() string {
	buf := make([]byte, 16)

	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(buf)
}

func (d *db) GetLock(key string, seconds int) (string, error) {
	value := randomValue()
	if err := d.SetNotExists(key, value, seconds); err != nil {
		return "", err
	}
	return value, nil
}

const setLockTTLScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("EXPIRE",KEYS[1], ARGV[2])
else
return 0
end`

func (d *db) SetLockTTL(key string, lockID string, seconds int) error {
	scr := redis.NewScript(1, setLockTTLScript)
	conn := d.pool.Get()
	defer conn.Close()

	lockIDData, _ := json.Marshal(lockID)

	ret, err := redis.Int(scr.Do(conn, key, lockIDData, seconds))
	if err != nil {
		return err
	}

	if ret == 0 {
		return errors.New("key does not exist or has expired")
	}

	return nil
}

func (l *db) DelLock(key string, lockID string) error {
	return l.DelKeyForValue(key, lockID)
}

func (l *db) Pool() *redis.Pool {
	return l.pool
}
