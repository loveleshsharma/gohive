package gohive

import (
	"errors"
	"sync/atomic"
)

type PoolState int

const (
	OPEN PoolState = iota
	CLOSED
)

type Pool struct {
	poolChan         chan Runnable
	quitChan         chan bool
	state            PoolState
	size             int
	availableWorkers int32
}

func NewFixedPool(size int) *Pool {
	pool := &Pool{
		poolChan: make(chan Runnable),
		state:    OPEN,
		size:     size,
	}

	for i := 0; i < size; i++ {
		go pool.worker()
		pool.availableWorkers++
	}

	return pool
}

func (p *Pool) Close() error {
	if p.state == OPEN {
		p.state = CLOSED
		return nil
	}
	return errors.New("error: cannot close an already closed pool")
}

func (p *Pool) Submit(r Runnable) error {
	if r == nil {
		return errors.New("cannot submit nil Runnable")
	}

	if atomic.LoadInt32(&p.availableWorkers) == 0 {
		return errors.New("cannot submit, pool is full")
	}

	p.poolChan <- r
	return nil
}

func (p *Pool) worker() {
	for {
		select {
		case r := <-p.poolChan:
			atomic.AddInt32(&p.availableWorkers, -1)
			r.Run()
			atomic.AddInt32(&p.availableWorkers, 1)
		case <-p.quitChan:
			break
		}
	}
}
