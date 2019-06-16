package GoHive

import (
	"fmt"
	"github.com/pkg/errors"
)

const DEFAULT_POOL_CAPACITY = 3

var (
	ErrInvalidPoolSize = errors.New("Invalid Pool Size: Pool Size must be a positive number!")
)

type RoutineService struct {
	routinePool  Pool
	waitingQueue *WaitingQueue
	poolSize     int
}

func NewDefaultRoutinePool() *RoutineService {
	wtQueue := NewWaitingQueue()
	routineService := RoutineService{waitingQueue: &wtQueue, poolSize: DEFAULT_POOL_CAPACITY}
	pool := NewFixedSizePool(DEFAULT_POOL_CAPACITY, &routineService)
	routineService.routinePool = pool
	return &routineService
}

func NewFixedSizeRoutinePool(numOfRoutines int) (*RoutineService, error) {

	if numOfRoutines <= 0 {
		panic(ErrInvalidPoolSize)
	}

	wtQueue := NewWaitingQueue()
	routineService := RoutineService{waitingQueue: &wtQueue, poolSize: numOfRoutines}
	pool := NewFixedSizePool(numOfRoutines, &routineService)
	routineService.routinePool = pool
	return &routineService, nil
}

func (rs *RoutineService) Submit(fun func()) error {
	newTask := Task{executable: fun}

	//	if worker is available, immediately assigning the task
	if rs.routinePool.isWorkerAvailable() {
		fmt.Println("Assigning!")
		err := rs.routinePool.assignTask(newTask)
		if err != nil {
			panic(err)
		}
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
	return rs.routinePool.capacity - rs.routinePool.runningWorkers
}

func (rs *RoutineService) Close() {
	rs.routinePool.Close()
}
