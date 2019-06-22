package GoHive

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	DefaultPoolCapacity = 10
)

var (
	ErrInvalidPoolSize = errors.New("Invalid pool Size: pool Size must be a positive number!")

	ErrInvalidPoolState = errors.New("pool is Closed: Cannot Assign task to a closed pool!")
)

type RoutineService struct {
	routinePool  *pool
	waitingQueue *WaitingQueue
	poolSize     int
}

func NewDefaultRoutinePool() *RoutineService {

	routineService := &RoutineService{
		waitingQueue: NewWaitingQueue(),
		poolSize:     DefaultPoolCapacity,
	}
	pool := newFixedSizePool(DefaultPoolCapacity, routineService)
	routineService.routinePool = pool

	return routineService
}

func NewFixedSizeRoutinePool(numOfRoutines int) (*RoutineService, error) {

	if numOfRoutines <= 0 {
		panic(ErrInvalidPoolSize)
	}
	routineService := &RoutineService{
		waitingQueue: NewWaitingQueue(),
		poolSize:     numOfRoutines,
	}
	pool := newFixedSizePool(numOfRoutines, routineService)
	routineService.routinePool = pool

	return routineService, nil
}

func (rs *RoutineService) Submit(fun func()) error {
	if fun == nil {
		return nil
	}

	if rs.routinePool.status == CLOSED {
		panic(ErrInvalidPoolState)
	}

	newTask := Task{executable: fun}

	//	if worker is available, immediately assigning the task
	if rs.routinePool.isWorkerAvailable() {
		fmt.Println("Assigning!")
		rs.routinePool.assignTask(newTask)
	} else {
		fmt.Println("Queuing!")
		rs.waitingQueue.EnqueueTask(newTask)
	}
	return nil
}

//notifies the routineService that one of the worker from
//the pool has completed its task and a new task can be
//assigned to this worker if waiting in the queue.
func (rs *RoutineService) notify() {
	if rs.waitingQueue.IsNotEmpty() {
		task, err := rs.waitingQueue.DequeueTask()
		if err != nil {
			fmt.Errorf("Error Dequeueing Task!")
			return
		}
		fmt.Println("Pulling task from Queue!")
		rs.routinePool.assignTask(task)
	}
}

func (rs *RoutineService) RunningWorkers() int {
	return rs.routinePool.runningWorkers
}

func (rs *RoutineService) PoolCapacity() int {
	return rs.routinePool.capacity
}

func (rs *RoutineService) AvailableWorkers() int {
	return rs.routinePool.capacity - rs.routinePool.runningWorkers //TODO: check for atomicity
}

func (rs *RoutineService) Close() {
	rs.routinePool.close()
}
