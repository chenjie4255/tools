package redis

import "github.com/gomodule/redigo/redis"

type FetchFunc func() (interface{}, error)

type DB interface {
	Get(key string, v interface{}) error
	Set(key string, value interface{}, expires int) error
	SetNotExists(key string, value interface{}, expires int) error
	Del(key string) error

	// DelByKeys deprecated Use DelKeysByScan instead
	DelByKeys(keyPattern string) error

	DelKeysByScan(keyPattern string) error
	DelKeys(keys []string) error
	DelKeyForValue(key string, value interface{}) error
	DelKeyForBytes(key string, value []byte) error
	ReplaceValue(key string, oldValid, newValue interface{}) error

	CmpAndSet(cmpKey string, cmpValue interface{}, setKey string, setValue interface{}) error

	CmpGTDecr(cmpKey string, greatThan int64) (int64, error)

	GetLock(key string, seconds int) (string, error)
	DelLock(key string, lockID string) error
	SetLockTTL(key string, lockID string, ttl int) error
	IncrByUint64(key string, step uint64) (uint64, error)
	IncrToUint64(key string, val uint64) (uint64, error)
	DecrByUint64(key string, step uint64) (uint64, error)

	Exists(key string) (bool, error)

	//
	CacheGet(key string, v interface{}, fn FetchFunc, expires int) error
	CacheGetB(key string, v interface{}, fn FetchFunc, expires int) (bool, error)

	AddSortSetStr(key string, value string, sortKey int64) error
	GetSortSetCount(key string, sortKeyFrom, sortKeyTo int64) (int, error)
	GetSortSetRangeStr(key string, sortKeyFrom, sortKeyTo int64) ([]string, error)
	RemoveSortSet(key string, sortKeyFrom, sortKeyTo int64) (int, error)

	PushStringList(key string, value string, expires int) error
	GetStringList(key string) ([]string, error)

	TTL(key string) (int, error)
	Time() (int64, error)

	SAdd(key string, values ...[]byte) error
	SRandMember(key string, count int) ([][]byte, error)
	SCard(key string) (int, error)

	Subscribe(channels []string, done chan bool, recv func(name string, data []byte)) error
	Publish(channel string, data []byte) (int, error)

	Pool() *redis.Pool
}

type WatchCBFn func(delete bool, data []byte)

type ShareInfo interface {
	Set(val string) error
	Get() (string, error)
	Del(oldVal string) error
}
