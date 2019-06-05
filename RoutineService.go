package GoHive

import "fmt"

const DEFAULT_POOL_CAPACITY = 3

type RoutineService struct {
	routinePool  Pool
	waitingQueue *WaitingQueue
	poolSize     int
}

func NewDefaultRoutinePool() RoutineService {
	wtQueue := NewWaitingQueue()
	routineService := RoutineService{waitingQueue: &wtQueue, poolSize: DEFAULT_POOL_CAPACITY}
	pool := NewFixedSizePool(DEFAULT_POOL_CAPACITY,&routineService)
	routineService.routinePool = pool
	return routineService
}

func NewFixedSizeRoutinePool(numOfRoutines int) RoutineService {
	wtQueue := NewWaitingQueue()
	routineService := RoutineService{waitingQueue: &wtQueue, poolSize: numOfRoutines}
	pool := NewFixedSizePool(numOfRoutines,&routineService)
	routineService.routinePool = pool
	return routineService
}

func (rs *RoutineService) Submit(fun func()) {
	newTask := Task{executable: fun}

	//	if worker is available, immediately assigning the task
	if rs.routinePool.isWorkerAvailable() {
		fmt.Println("Assigning!")
		rs.routinePool.assignTask(newTask)
	} else {
		fmt.Println("Queuing!")
		rs.waitingQueue.EnqueueTask(newTask)
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
