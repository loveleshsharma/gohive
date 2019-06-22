package GoHive

import (
	"sync"
)

type state int

const(
	OPEN   state = 1
	CLOSED state = 0
)

type pool struct {
	capacity         int
	runningWorkers   int
	availableWorkers sync.Pool
	closePool		 sync.Once
	locker			 sync.Mutex
	status           state
	routineService   *RoutineService //reference back to the routine service who owns this pool
}

func newFixedSizePool(newSize int, routineService *RoutineService) *pool {
	newPool := pool{
		capacity: newSize,
		routineService: routineService,
		status: OPEN,
		availableWorkers: sync.Pool{
			New: func() interface{} {
				return new(worker)
			},
		},
	}
	return &newPool
}

func (p *pool) assignTask(task Task) {
	p.locker.Lock()
	defer p.locker.Unlock()
	worker := p.availableWorkers.Get().(*worker)
	worker.taskChan = make(chan func())
	worker.pool = p
	go worker.run()
	worker.taskChan <- task.getTask()
	p.runningWorkers++
}

func (p *pool) done(w *worker) {
	w = new(worker)
	p.availableWorkers.Put(w)
	p.locker.Lock()
	p.runningWorkers--
	p.locker.Unlock()
	p.routineService.notify() //	notify the RoutineService to pull more tasks from the queue if waiting
}

func (p *pool) close() {
	p.closePool.Do(func() {
		p.status = CLOSED
	})
}

func (p *pool) isWorkerAvailable() bool {
	return p.capacity > p.runningWorkers
}

func (p *pool) getRunning() int {
	return p.runningWorkers
}

func (p *pool) getCapacity() int {
	return p.capacity
}
