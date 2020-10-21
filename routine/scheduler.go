package routine

import (
	"fmt"
	"github.com/chenjie4255/errors"
	"github.com/chenjie4255/tools/errcode"
	"github.com/chenjie4255/tools/log"
	"github.com/chenjie4255/tools/redis"
	"time"
)

var logger *log.Logger

func init() {
	logger = log.NewLoggerWithSentry("routine")
}

type Scheduler interface {
	FetchJobs() ([]Job, bool, error)
	SetInitialPos(pos SchedulerPos)
	SetInitialPosWithServerTime()
	LastPos() SchedulerPos
}

type SchedulerConfig struct {
	PartitionSteps uint32
	AheadSecond    uint32
	DelaySecond    uint32
}

type scheduler struct {
	redisDB redis.DB
	config  SchedulerConfig
	tb      WeeklyTable

	redisKeySchedulePos   string
	redisKeyLastResetTime string

	lastPos SchedulerPos
}

type defTimeGetter struct{}

func (g defTimeGetter) Now() time.Time {
	return time.Now()
}

func NewScheduler(tb WeeklyTable, redisDB redis.DB, config SchedulerConfig) Scheduler {
	ret := scheduler{}
	ret.tb = tb
	ret.redisDB = redisDB
	ret.config = config

	ret.redisKeySchedulePos = fmt.Sprintf("schedule_pos_%s", tb.UniqueName())
	ret.redisKeyLastResetTime = fmt.Sprintf("last_reset_time_%s", tb.UniqueName())

	return &ret
}

func (s *scheduler) SetInitialPos(pos SchedulerPos) {
	s.redisDB.SetNotExists(s.redisKeySchedulePos, pos, 0)
}

func (s *scheduler) SetInitialPosWithServerTime() {
	t, _ := s.redisDB.Time()
	s.SetInitialPos(newSchedulerPos(t, 0))
}

func curTimeOffset() uint32 {
	tNow := time.Now()
	return getWeekOffsetFromTime(tNow)
}

func getWeekOffsetFromTime(t time.Time) uint32 {
	utcT := t.UTC()
	wd := utcT.Weekday()
	return uint32(int(wd)*3600*24 + utcT.Hour()*3600 + utcT.Minute()*60 + utcT.Second())
}

func addWeekOffset(base uint32, delta int32) WeekOffset {
	ret := int32(base) + delta
	if ret >= maxWeeklyOffset {
		return WeekOffset(ret - maxWeeklyOffset)
	}

	if ret < 0 {
		return WeekOffset(ret + maxWeeklyOffset)
	}

	return WeekOffset(ret)
}

func (s *scheduler) LastPos() SchedulerPos {
	return s.lastPos
}

func (s *scheduler) FetchJobs() ([]Job, bool, error) {
	var pos SchedulerPos
	if val, err := s.redisDB.IncrByUint64(s.redisKeySchedulePos, uint64(s.config.PartitionSteps)); err != nil {
		logger.AddFile().WithError(err).Error("failed to incr schedule pos")
		return nil, false, err
	} else {
		pos = SchedulerPos(val)
		s.lastPos = pos
	}

	if s.needRest(pos) {
		return nil, true, nil
	}

	if s.needReset(pos) {
		tNow, err := s.redisDB.Time()
		if err != nil {
			logger.AddFile().WithError(err).Error("failed to get redis server time")
			return nil, false, err
		}

		expireTime := s.config.DelaySecond
		if expireTime == 0 {
			expireTime = 10
		}
		if err := s.redisDB.SetNotExists(s.redisKeyLastResetTime, tNow, int(expireTime)); err != nil {
			if errors.FindTag(err, errcode.ResExisted) {
				return nil, false, nil
			}

			logger.AddFile().WithError(err).Error("failed to reset schedule pos(set last reset time step)")
		}

		newPos := newSchedulerPos(tNow, 0)
		if err := s.redisDB.Set(s.redisKeySchedulePos, newPos, 0); err != nil {
			logger.AddFile().WithError(err).Error("failed to reset schedule pos")
			return nil, false, err
		}

		logger.AddFile().Info("reset schedule's pos successfully")
		return nil, false, nil
	}

	partition := pos.partition()
	offset := pos.offset()

	fromPartition := partition - s.config.PartitionSteps

	if fromPartition >= s.tb.PartitionCount() {
		// next second
		ts := pos.timestamp() + 1 // next second
		newPos := newSchedulerPos(ts, 0)
		if val, err := s.redisDB.IncrToUint64(s.redisKeySchedulePos, uint64(newPos)); err != nil {
			logger.AddFile().WithError(err).Error("failed to move schedule pos")
			return nil, false, err
		} else if val > 0 {
			//fmt.Printf("succeed incr pos by 1 second, now: %d(%d)\n", val, ts)
		}

		return nil, false, nil
	}

	ret, err := s.tb.ScanCellsPartitions(uint64(offset), uint64(fromPartition), uint64(partition))
	if err != nil {
		return nil, false, err
	}
	jobs := ret.Jobs()
	//for i := range jobs {
	//}
	//fmt.Printf("[%d]query: (%d) [%d - %d] --> %d\n", pos.timestamp(), offset, fromPartition, partition, len(jobs))

	return jobs, false, nil
}

func (s *scheduler) needRest(pos SchedulerPos) bool {
	ts := pos.timestamp()
	tNow := time.Now().Unix()

	return ts-tNow > int64(s.config.AheadSecond) && pos.Partition() > s.tb.PartitionCount()
}

func (s *scheduler) needReset(pos SchedulerPos) bool {
	ts := pos.timestamp()
	tNow := time.Now().Unix()

	if tNow-ts > int64(s.config.DelaySecond) {
		return true
	}

	return false
}
