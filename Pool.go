package GoHive

import (
	"github.com/pkg/errors"
)

type state int

const (
	OPEN   state = 1
	CLOSED state = 0
)

var (
	ErrInvalidPoolState = errors.New("Pool is Closed: Cannot Assign task to a closed pool!")
)

type Pool struct {
	capacity       int
	runningWorkers int
	workers        []Worker
	status         state
	routineService *RoutineService //reference back to the routine service who owns this pool
}

func NewFixedSizePool(newSize int, routineService *RoutineService) Pool {
	newPool := Pool{capacity: newSize, workers: make([]Worker, newSize), routineService: routineService, status: OPEN}
	return newPool
}

func (p *Pool) assignTask(task Task) error {
	if p.status == OPEN {
		for i := range p.workers {
			if p.workers[i].isOccupied == false {
				p.workers[i] = Worker{taskChan: make(chan func()), pool: p, isOccupied: true}
				go p.workers[i].run()
				p.workers[i].taskChan <- task.getTask()
				p.runningWorkers++
				break
			}
		}
		return nil
	}
	return ErrInvalidPoolState
}

func (p *Pool) Done(w *Worker) {
	p.runningWorkers--
	w.isOccupied = false
	p.routineService.notify() //	notify the RoutineService to pull more tasks from the queue if waiting
}

func (p *Pool) Close() {
	for i := range p.workers {
		if p.workers[i].isOccupied == false {
			p.workers[i] = Worker{}
		}
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
