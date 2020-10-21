package worker

import (
	"sync/atomic"
	"testing"
)

func TestWorker(t *testing.T) {
	pool := NewPool(16)

	var sum int32
	for i := 0; i < 10; i++ {
		val := i
		pool.AddJob(func() {
			atomic.AddInt32(&sum, int32(val))
			pool.JobDone()
		})
	}

	pool.WaitAll()

	if sum != 45 {
		t.Errorf("result error, should equal 45 but got:%d", sum)
	}
}

func TestAddJobWait(t *testing.T) {
	pool := NewPool(16)

	var sum int32
	for i := 0; i < 10; i++ {
		val := i
		pool.AddJobWait(func() {
			atomic.AddInt32(&sum, int32(val))
		})
	}

	pool.WaitAll()

	if sum != 45 {
		t.Errorf("result error, should equal 45 but got:%d", sum)
	}
}
