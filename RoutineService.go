package GoHive

import "fmt"

const DEFAULT_POOL_CAPACITY = 10

type RoutineService struct {
	routinePool  Pool
	waitingQueue *WaitingQueue
	poolSize     int
}

func NewDefaultRoutinePool() RoutineService {
	wtQueue := NewWaitingQueue()
	return RoutineService{waitingQueue: &wtQueue, poolSize: DEFAULT_POOL_CAPACITY, routinePool: NewFixedSizePool(DEFAULT_POOL_CAPACITY)}
}

func NewFixedSizeRoutinePool(numOfRoutines int) RoutineService {
	wtQueue := NewWaitingQueue()
	return RoutineService{waitingQueue: &wtQueue, poolSize: numOfRoutines}
}

func (rs *RoutineService) Submit(fun func()) {
	newTask := Task{executable: fun}
	rs.waitingQueue.EnqueueTask(newTask)

	//	checking availability in the pool
	if rs.routinePool.isWorkerAvailable() {
		newTask, err := rs.waitingQueue.DequeueTask()
		if err != nil {
			fmt.Println("Cannot Dequeue task:", err.Error())
		}
		rs.routinePool.assignTask(newTask)
	}
}

func (rs *RoutineService) RunningWorkers() int {
	return rs.routinePool.running
}

func (rs *RoutineService) PoolCapacity() int {
	return rs.routinePool.capacity
}

func (rs *RoutineService) AvailableWorkers() int {
	return rs.routinePool.capacity - rs.routinePool.running
}
