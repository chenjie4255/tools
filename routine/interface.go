package routine

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func cellIndex(offset, piece uint32) uint64 {
	return uint64(offset)<<32 | uint64(piece)
}

// CellIndex , high 32bit: offset, low 32bit: partition number
type CellIndex uint64

type Cell struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	Offset    int64         `json:"offset" bson:"offset"`
	Partition int64         `json:"partition" bson:"partition"`
	Index     string        `json:"-" bson:"index"`
	Jobs      []Job         `json:"jobs" bson:"jobs"`
}

type Job struct {
	UID  string `json:"uid" bson:"uid"`
	Data []byte `json:"data" bson:"data"`
}

func newSchedulerPos(ts int64, partition int64) SchedulerPos {
	return SchedulerPos(ts<<32 | partition)
}

// SchedulerPos, high 32bit: timestamp, low 32bit: partition number
type SchedulerPos uint64

const (
	maxWeeklyOffset = 3600 * 24 * 7
)

func (p SchedulerPos) timestamp() int64 {
	return int64(p >> 32)
}

func (p SchedulerPos) offset() uint32 {
	ts := int64(p >> 32)
	t := time.Unix(ts, 0).UTC()
	return getWeekOffsetFromTime(t)
}

func (p SchedulerPos) WeekOffset() uint32 {
	return p.offset()
}

func (p SchedulerPos) Partition() uint32 {
	return uint32(p & 0xffffffff)
}

func (p SchedulerPos) partition() uint32 {
	return uint32(p & 0xffffffff)
}

type WeekOffset uint32

type WeeklyTableIndex struct {
	Offset    uint32
	Partition uint32
}

type Cells []Cell

func (c Cells) Jobs() []Job {
	ret := []Job{}
	for i := range c {
		ret = append(ret, c[i].Jobs...)
	}

	return ret
}

func NowWeekOffset() uint32 {
	tNow := time.Now().UTC()
	return getWeekOffsetFromTime(tNow)
}

type WeeklyTable interface {
	AddJob(job Job, offsets []WeekOffset) ([]string, error)

	// ScanCellsPartitions scan cells within [from, to)
	ScanCellsPartitions(offset, from, to uint64) (Cells, error)

	// Remove
	RemoveJob(uid string, indexes []string) error

	PartitionCount() uint32

	UniqueName() string
}
