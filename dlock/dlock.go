package dlock

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/chenjie4255/tools/log"
	rediscm "github.com/chenjie4255/tools/redis"
)

var ErrorTimeout = errors.New("operation timeout")

type DLock interface {
	GetLock(key string, seconds int) (string, error)
	GetLockWait(key string, seconds int, timeout time.Duration) (string, error)
	DelLock(key string, secret string) error
	SetExpiredTime(key, secret string, seconds int) error
}

var logger *log.Logger

func init() {
	logger = log.NewLoggerWithSentry("dlock")
}

func New(host, password string, dbNum int) DLock {
	pool := rediscm.NewPool(host, password, dbNum)
	return &dlock{pool}
}

func NewWithPool(pool *redis.Pool) DLock {
	return &dlock{pool}
}

type dlock struct {
	pool *redis.Pool
}

func randomValue() string {
	buf := make([]byte, 16)

	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(buf)
}

const script = `if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("del",KEYS[1])
else
return 0
end`

func (l *dlock) GetLock(key string, seconds int) (string, error) {
	conn := l.pool.Get()
	defer conn.Close()

	value := randomValue()

	if _, err := redis.String(conn.Do("SET", key, value, "EX", seconds, "NX")); err != nil {
		return "", err
	}

	return value, nil
}

func (l *dlock) GetLockWait(key string, seconds int, timeout time.Duration) (string, error) {
	conn := l.pool.Get()
	defer conn.Close()

	value := randomValue()

	errChan := make(chan error, 1)

	go func() {
		for {
			if _, err := redis.String(conn.Do("SET", key, value, "EX", seconds, "NX")); err != nil {
				if err == redis.ErrNil {
					time.Sleep(timeout / 10)
					continue
				} else {
					errChan <- err
					return
				}
			} else {
				errChan <- nil
				return
			}
		}
	}()

	select {
	case <-time.After(timeout):
		return "", ErrorTimeout
	case err := <-errChan:
		if err != nil {
			return "", err
		}
		return value, nil
	}
}

const setExpireScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
return redis.call("EXPIRE",KEYS[1], ARGV[2])
else
return 0
end`

func (l *dlock) SetExpiredTime(key string, secret string, seconds int) error {
	// conn := l.pool.Get()
	scr := redis.NewScript(1, setExpireScript)
	conn := l.pool.Get()
	defer conn.Close()

	ret, err := redis.Int(scr.Do(conn, key, secret, seconds))
	if err != nil {
		return err
	}

	if ret == 0 {
		return errors.New("key does not exist or has expired")
	}

	return nil
}

func (l *dlock) DelLock(key string, secret string) error {
	scr := redis.NewScript(1, script)
	conn := l.pool.Get()
	defer conn.Close()

	ret, err := redis.Int(scr.Do(conn, key, secret))
	if err != nil {
		return err
	}

	if ret == 0 {
		return errors.New("key does not exist or has expired")
	}

	return nil

}
