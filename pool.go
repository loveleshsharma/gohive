package gohive

import (
	"sync"
)

type state int

const (
	//OPEN indicates that the pool is open to accept new tasks
	OPEN   state = 1

	//CLOSED indicates that the pool is closed and won't accept new tasks
	CLOSED state = 0
)

//pool represents a group of workers to whom tasks can be assigned.
type pool struct {

	//number of workers in the pool
	poolCapacity int

	//number of currently active workers
	activeWorkers int

	//pool of available workers out of total poolCapacity
	availableWorkers sync.Pool

	//object which closes the pool and it can be called only once in the program scope
	closePool sync.Once

	//Mutex used for atomic operations
	locker sync.Mutex

	//represents the current state of the pool(OPEN/CLOSED)
	status state

	//reference back to the routine service who owns this pool
	poolService *PoolService
}

//returns an instance of pool with the size specified
func newPool(newSize int, poolService *PoolService) *pool {
	newPool := pool{
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

//gets an available worker from the pool and assigns the task
func (p *pool) assignTask(task Task) {
	p.locker.Lock()
	defer p.locker.Unlock()
	worker := p.availableWorkers.Get().(*worker)
	worker.taskChan = make(chan func())
	worker.pool = p
	worker.run()
	worker.taskChan <- task.getTask()
	p.activeWorkers++
}

//done is called by a worker after completing its task in order to
//notify the poolService that a worker is now available in the pool
//and waiting tasks can be pulled from the Queue if any
func (p *pool) done(w *worker) {
	w = new(worker)
	p.availableWorkers.Put(w)
	p.locker.Lock()
	p.activeWorkers--
	p.locker.Unlock()
	p.poolService.notify() //	notify the PoolService to pull more tasks from the queue if waiting
}

//closes the pool and makes sure
//that no more tasks should be accepted
func (p *pool) close() {
	p.closePool.Do(func() {
		p.status = CLOSED
	})
}

//returns if any worker is available
func (p *pool) isWorkerAvailable() bool { return p.poolCapacity > p.activeWorkers }

//returns number of active workers
func (p *pool) getActiveWorkers() int { return p.activeWorkers }

//returns the pool size
func (p *pool) getPoolSize() int { return p.poolCapacity }
