package GoHive

import (
	"github.com/pkg/errors"
	"sync"
)

type state int

const (
	OPEN   state = 1
	CLOSED state = 0
)

var (
	ErrInvalidPoolState = errors.New("Pool is Closed: Cannot Assign task to a closed pool!")
)

//TODO: implement lock mechanishm in Pool for atomic operation at a time

type Pool struct {
	capacity         int
	runningWorkers   int
	availableWorkers sync.Pool
	status           state
	routineService   *RoutineService //reference back to the routine service who owns this pool
}

func NewFixedSizePool(newSize int, routineService *RoutineService) Pool {
	newPool := Pool{capacity: newSize, routineService: routineService, status: OPEN}
	newPool.availableWorkers = sync.Pool{
		New: func() interface{} {
			return new(Worker)
		},
	}
	return newPool
}

func (p *Pool) assignTask(task Task) error {
	if p.status == OPEN {
		worker := p.availableWorkers.Get().(*Worker)
		worker.taskChan = make(chan func())
		worker.pool = p
		go worker.run()
		worker.taskChan <- task.getTask()
		p.runningWorkers++
		return nil
	}
	return ErrInvalidPoolState
}

func (p *Pool) Done(w *Worker) {
	w = new(Worker)
	p.availableWorkers.Put(w)
	p.runningWorkers--
	p.routineService.notify() //	notify the RoutineService to pull more tasks from the queue if waiting
}

func (p *Pool) Close() {
	if p.status == CLOSED {
		return
	}
	p.status = CLOSED
}

func (p *Pool) isWorkerAvailable() bool {
	return p.capacity > p.runningWorkers
}

func (p *Pool) getRunning() int {
	return p.runningWorkers
}

func (p *Pool) getCapacity() int {
	return p.capacity
}
