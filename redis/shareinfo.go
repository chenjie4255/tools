package redis

import (
	"fmt"
	"github.com/chenjie4255/tools/errcode"
	"github.com/chenjie4255/tools/gor"
	"github.com/chenjie4255/errors"
	"sync"
	"time"
)

const (
	ShareInfoUpdatedEvent = "UPDATED"
	ShareInfoDeletedEvent = "DELETED"
)

type shareInfo struct {
	redisDB DB

	rwlock sync.RWMutex

	// keys
	channelKey string
	valueKey   string

	data             string
	refreshInterval  int
	refreshCountDown int
}

func NewShareInfo(infoName string, redisDB DB, forceRefresh int) ShareInfo {
	ret := shareInfo{}
	ret.redisDB = redisDB

	ret.valueKey = fmt.Sprintf("redis_share_info_%s", infoName)
	ret.channelKey = fmt.Sprintf("redis_share_info_channel_%s", infoName)
	ret.refreshInterval = forceRefresh
	ret.refreshCountDown = forceRefresh

	ret.data, _ = ret.getFromRedis()

	done := make(chan bool)
	go ret.asyncWatch(done)
	<-done

	ret.loopRefresh()

	return &ret
}

func (i *shareInfo) loopRefresh() {
	gor.RunWithRecover(func() {
		for {
			time.Sleep(1 * time.Second)
			i.refreshCountDown--
			if i.refreshCountDown <= 0 {
				i.data, _ = i.getFromRedis()
				i.refreshCountDown = i.refreshInterval
			}
		}
	})
}

func (i *shareInfo) asyncWatch(done chan bool) {
	i.redisDB.Subscribe([]string{i.channelKey}, done, func(channel string, data []byte) {
		if string(data) == ShareInfoUpdatedEvent {
			val, err := i.getFromRedis()
			if err == nil {
				i.rwlock.Lock()
				i.data = val
				i.rwlock.Unlock()
				i.refreshCountDown = i.refreshInterval
			}
		} else if string(data) == ShareInfoDeletedEvent {
			i.rwlock.Lock()
			i.data = ""
			i.rwlock.Unlock()
		}
	})
}

func (i *shareInfo) Set(val string) error {
	if err := i.redisDB.Set(i.valueKey, val, 0); err != nil {
		return err
	}

	i.data = val

	_, err := i.redisDB.Publish(i.channelKey, []byte(ShareInfoUpdatedEvent))
	return err
}

func (i *shareInfo) getFromRedis() (string, error) {
	ret := ""
	err := i.redisDB.Get(i.valueKey, &ret)
	if errors.FindTag(err, errcode.ResNotFound) {
		return "", nil
	}
	return ret, err
}

func (i *shareInfo) Get() (string, error) {
	if i.data != "" {
		return i.data, nil
	}

	return i.getFromRedis()
}

func (i *shareInfo) Del(oldVal string) error {
	if err := i.redisDB.DelKeyForValue(i.valueKey, oldVal); err != nil {
		return err
	}

	i.data = ""

	_, err := i.redisDB.Publish(i.channelKey, []byte(ShareInfoDeletedEvent))
	return err
}
