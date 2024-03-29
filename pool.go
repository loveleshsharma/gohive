package gohive

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const (
	OPEN = iota
	CLOSED
)

type Pool struct {
	poolChan         chan Runner
	quitChan         chan bool
	state            int32
	size             int
	availableWorkers int32
}

func NewFixedPool(size int) *Pool {
	pool := &Pool{
		poolChan: make(chan Runner),
		quitChan: make(chan bool),
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
	if atomic.CompareAndSwapInt32(&p.state, OPEN, CLOSED) {
		for i := 0; i < p.size; i++ {
			p.quitChan <- true
		}
		fmt.Println("pool is closed.")
		return nil
	}
	return errors.New("error: cannot close an already closed pool")
}

func (p *Pool) IsPoolClosed() bool {
	return atomic.LoadInt32(&p.state) == CLOSED
}

func (p *Pool) Submit(r Runner) error {
	if r == nil {
		return errors.New("cannot submit nil Runner")
	}

	if atomic.LoadInt32(&p.state) == CLOSED {
		return errors.New("cannot submit, pool is closed")
	}

	p.poolChan <- r
	return nil
}

func (p *Pool) worker() {
	defer fmt.Println("closing...")
Loop:
	for {
		select {
		case r := <-p.poolChan:
			atomic.AddInt32(&p.availableWorkers, -1)
			r.Run()
			atomic.AddInt32(&p.availableWorkers, 1)
		case <-p.quitChan:
			break Loop
		}
	}
}
