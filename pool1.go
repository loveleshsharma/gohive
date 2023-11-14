package gohive

import (
	"sync"
)

type state int

const (
	//OPEN indicates that the pool1 is open to accept new tasks
	OPEN state = 1

	//CLOSED indicates that the pool1 is closed and won't accept new tasks
	CLOSED state = 0
)

// pool1 represents a group of workers to whom tasks can be assigned.
type pool1 struct {
	//number of workers in the pool1
	poolCapacity int

	//number of currently active workers
	activeWorkers int

	//pool1 of available workers out of total poolCapacity
	availableWorkers sync.Pool

	//object which closes the pool1 and it can be called only once in the program scope
	closePool sync.Once

	//Mutex used for atomic operations
	locker sync.Mutex

	//represents the current state of the pool1(OPEN/CLOSED)
	status state

	//reference back to the routine service who owns this pool1
	poolService *PoolService
}

// returns an instance of pool1 with the size specified
func newPool(newSize int, poolService *PoolService) *pool1 {
	newPool := pool1{
		poolCapacity: newSize,
		poolService:  poolService,
		status:       OPEN,
		availableWorkers: sync.Pool{
			New: func() interface{} {
				return new(worker)
			},
		},
	}
	return &newPool
}

// gets an available worker from the pool1 and assigns the task
func (p *pool1) assignTask(task Task) {
	p.locker.Lock()
	defer p.locker.Unlock()
	worker := p.availableWorkers.Get().(*worker)
	worker.taskChan = make(chan func())
	worker.pool = p
	worker.run()
	worker.taskChan <- task.getTask()
	p.activeWorkers++
}

// done is called by a worker after completing its task in order to
// notify the poolService that a worker is now available in the pool1
// and waiting tasks can be pulled from the Queue if any
func (p *pool1) done(w *worker) {
	w = new(worker)
	p.availableWorkers.Put(w)
	p.locker.Lock()
	p.activeWorkers--
	p.locker.Unlock()
	p.poolService.notify() //	notify the PoolService to pull more tasks from the queue if waiting
}

// closes the pool1 and makes sure
// that no more tasks should be accepted
func (p *pool1) close() {
	p.closePool.Do(func() {
		p.status = CLOSED
	})
}

// returns if any worker is available
func (p *pool1) isWorkerAvailable() bool { return p.poolCapacity > p.activeWorkers }

// returns number of active workers
func (p *pool1) getActiveWorkers() int { return p.activeWorkers }

// returns the pool1 size
func (p *pool1) getPoolSize() int { return p.poolCapacity }
