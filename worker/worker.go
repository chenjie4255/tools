package worker

import (
	"sync"

	"github.com/ivpusic/grpool"
)

type Pool struct {
	pool   *grpool.Pool
	errors []error
	mux    sync.RWMutex
}

func NewPool(num int) *Pool {
	return &Pool{grpool.NewPool(num, num), nil, sync.RWMutex{}}
}

// AddJob deprecated Use AddJobWait instead
func (p *Pool) AddJob(job func()) {
	p.pool.WaitCount(1)
	p.pool.JobQueue <- job
}

func (p *Pool) AddJobWait(job func()) {
	p.pool.WaitCount(1)
	p.pool.JobQueue <- func() {
		defer p.JobDone()
		job()
	}
}

func (p *Pool) JobDone() {
	p.pool.JobDone()
}

func (p *Pool) WaitAll() {
	p.pool.WaitAll()
}

func (p *Pool) Release() {
	p.pool.Release()
}

func (p *Pool) SetError(err error) {
	p.mux.Lock()
	p.errors = append(p.errors, err)
	p.mux.Unlock()
}

func (p *Pool) AnyError() bool {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return len(p.errors) > 0
}

func (p *Pool) Errors() []error {
	p.mux.RLock()
	ret := make([]error, len(p.errors))
	copy(ret, p.errors)
	p.mux.RUnlock()
	return ret
}
