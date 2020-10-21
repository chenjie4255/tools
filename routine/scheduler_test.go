package routine

import (
	"context"
	"fmt"
	"github.com/chenjie4255/tools/mongohelper"
	"github.com/chenjie4255/tools/redis"
	"github.com/chenjie4255/tools/testenv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_currentTimeOffset(t *testing.T) {
	got := curTimeOffset()
	t.Log(got)

	if got > maxWeeklyOffset {
		t.Fatalf("invalid current time offset, %d", got)
	}
}

func TestScheduler(t *testing.T) {
	if testing.Short() {
		t.Skip("short test")
	}

	env := testenv.GetIntegratedTestEnv()
	if env.MongoHost == "" || env.RedisHost == "" {
		t.Skip("the environment configuration is not ready yet")
	}
	mgoSession, _ := mongohelper.NewClient(env.MongoHost, "", "", "admin")

	tb := NewWeeklyTable("test", "weekly_jobs", 1024, mgoSession)
	mgoSession.Database("test").Drop(context.Background())

	redis.FlushDB(env.RedisHost, env.RedisPassword, 0)
	redisDb := redis.NewDB(env.RedisHost, env.RedisPassword, 0)

	cf := SchedulerConfig{}
	cf.AheadSecond = 30
	cf.DelaySecond = 30
	cf.PartitionSteps = 128

	jb := Job{}
	jb.UID = "uid_123"
	jb.Data = []byte(`123123123`)

	ct := time.Now()
	t.Logf("current t: %d", ct.Unix())
	co := getWeekOffsetFromTime(ct)

	result, err := tb.AddJob(jb, []WeekOffset{addWeekOffset(co, 1),
		addWeekOffset(co, 5),
		addWeekOffset(co, 2),
		addWeekOffset(co, 3),
		addWeekOffset(co, 4),
		addWeekOffset(co, 50), // out
		addWeekOffset(co, 10),
		addWeekOffset(co, 20),
	})

	assert.NoError(t, err)
	t.Logf("job indexes:%+v", result)

	var fetchJobCount int32
	var restCount int32

	for i := 0; i < 50; i++ {
		go func() {
			sc := NewScheduler(tb, redisDb, cf)
			sc.SetInitialPos(newSchedulerPos(ct.Unix(), 0))

			for {
				jobs, needRest, err := sc.FetchJobs()
				assert.NoError(t, err)
				atomic.AddInt32(&fetchJobCount, int32(len(jobs)))
				if needRest {
					time.Sleep(100 * time.Millisecond)
					atomic.AddInt32(&restCount, 1)
				}
			}
		}()
	}
	time.Sleep(2 * time.Second)
	t.Logf("rest count:%d", restCount)
	assert.Equal(t, int32(7), fetchJobCount)
}

func TestSchedulerAutoReset(t *testing.T) {
	if testing.Short() {
		t.Skip("short test")
	}
	//123
	env := testenv.GetIntegratedTestEnv()
	if env.MongoHost == "" || env.RedisHost == "" {
		t.Skip("the environment configuration is not ready yet")
	}

	mgoSession, _ := mongohelper.NewClient(env.MongoHost, "", "", "admin")

	tb := NewWeeklyTable("test", "weekly_jobs", 1024, mgoSession)
	mgoSession.Database("test").Drop(context.Background())

	redis.FlushDB(env.RedisHost, env.RedisPassword, 0)
	redisDb := redis.NewDB(env.RedisHost, env.RedisPassword, 0)

	cf := SchedulerConfig{}
	cf.AheadSecond = 10
	cf.DelaySecond = 10
	cf.PartitionSteps = 128

	ct := time.Now()
	t.Logf("current t: %d", ct.Unix())
	co := getWeekOffsetFromTime(ct)

	addOffsetJob := func(offset int32) {
		j := Job{}
		j.Data = []byte(`111`)
		j.UID = fmt.Sprintf("uid_%d", int32(co)+offset)
		tb.AddJob(j, []WeekOffset{addWeekOffset(co, offset)})
	}

	addOffsetJob(1)
	addOffsetJob(2)
	addOffsetJob(3)
	addOffsetJob(4)
	addOffsetJob(5)
	addOffsetJob(50)
	addOffsetJob(20)
	addOffsetJob(-5)
	addOffsetJob(-6)
	addOffsetJob(-7)

	var fetchJobCount int32
	var restCount int32
	mux := sync.Mutex{}
	jobIndex := []string{}

	for i := 0; i < 50; i++ {
		go func() {
			sc := NewScheduler(tb, redisDb, cf)

			// set schedule's initial time as current_time - 20s to simulate the context of auto reset
			sc.SetInitialPos(newSchedulerPos(ct.Unix()-20, 0))

			for {
				jobs, needRest, err := sc.FetchJobs()
				for i := range jobs {
					mux.Lock()
					jobIndex = append(jobIndex, jobs[i].UID)
					mux.Unlock()
				}
				assert.NoError(t, err)
				atomic.AddInt32(&fetchJobCount, int32(len(jobs)))
				if needRest {
					time.Sleep(100 * time.Millisecond)
					atomic.AddInt32(&restCount, 1)
				}
			}
		}()
	}
	time.Sleep(2 * time.Second)
	t.Logf("rest count:%d", restCount)
	t.Logf("UIDs:%+v", jobIndex)
	assert.Equal(t, int32(5), fetchJobCount)
}

func TestSchedulerFromPass(t *testing.T) {
	if testing.Short() {
		t.Skip("short test")
	}

	env := testenv.GetIntegratedTestEnv()
	if env.MongoHost == "" || env.RedisHost == "" {
		t.Skip("the environment configuration is not ready yet")
	}

	mgoSession, _ := mongohelper.NewClient(env.MongoHost, "", "", "admin")

	tb := NewWeeklyTable("test", "weekly_jobs", 1024, mgoSession)
	mgoSession.Database("test").Drop(context.Background())

	redis.FlushDB(env.RedisHost, env.RedisPassword, 0)
	redisDb := redis.NewDB(env.RedisHost, env.RedisPassword, 0)

	cf := SchedulerConfig{}
	cf.AheadSecond = 10
	cf.DelaySecond = 20
	cf.PartitionSteps = 128

	ct := time.Now().UTC()
	t.Logf("current t: %d", ct.Unix())
	co := getWeekOffsetFromTime(ct)

	addOffsetJob := func(offset int32) {
		j := Job{}
		j.Data = []byte(`111`)
		j.UID = fmt.Sprintf("uid_%d", int32(co)+offset)
		tb.AddJob(j, []WeekOffset{addWeekOffset(co, offset)})
	}

	addOffsetJob(1)
	addOffsetJob(2)
	addOffsetJob(5)
	addOffsetJob(20)
	addOffsetJob(-5)
	addOffsetJob(-6)
	addOffsetJob(-7)

	var fetchJobCount int32
	var restCount int32
	mux := sync.Mutex{}
	jobIndex := []string{}

	for i := 0; i < 50; i++ {
		go func() {
			sc := NewScheduler(tb, redisDb, cf)

			// set schedule's initial time as current_time - 20s to simulate the context of auto reset
			sc.SetInitialPos(newSchedulerPos(ct.Unix()-10, 0))

			for {
				jobs, needRest, err := sc.FetchJobs()
				for i := range jobs {
					mux.Lock()
					jobIndex = append(jobIndex, jobs[i].UID)
					mux.Unlock()
				}
				assert.NoError(t, err)
				atomic.AddInt32(&fetchJobCount, int32(len(jobs)))
				if needRest {
					time.Sleep(100 * time.Millisecond)
					atomic.AddInt32(&restCount, 1)
				}
			}
		}()
	}
	time.Sleep(2 * time.Second)
	t.Logf("rest count:%d", restCount)
	t.Logf("UIDs:%+v", jobIndex)
	assert.Equal(t, int32(6), fetchJobCount)
}
