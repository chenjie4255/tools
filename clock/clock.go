package clock

import "time"

type Clock interface {
	GetUnix() int64
}

type FixedClock struct {
	ts int64
}

func (c FixedClock) GetUnix() int64 {
	return c.ts
}

func (c FixedClock) SetTs(ts int64) {
	c.ts = ts
}

type OffsetClock struct {
	offset int64
}

func (c OffsetClock) GetUnix() int64 {
	return time.Now().Unix() + c.offset
}

func (c OffsetClock) SetOffset(offset int64) {
	c.offset = offset
}

func NewFixedClock(ts int64) *FixedClock {
	return &FixedClock{ts}
}

func NewOffsetClock(offset int64) *OffsetClock {
	return &OffsetClock{offset}
}
