package GoHive

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

//default pool size for number of workers in the pool
const DefaultPoolSize = 10

//error for invalid pool size
var ErrInvalidPoolSize = errors.New("Invalid pool Size: pool Size must be a positive number!")

//error for invalid pool state
var	ErrInvalidPoolState = errors.New("pool is Closed: Cannot Assign task to a closed pool!")


//PoolService acts as an orchestrator of the entire GoHive functionality
//It consists of a pool, that manages workers that run tasks and it
//consists of a taskQueue that holds tasks waiting to acquire a worker
type PoolService struct {

	//pool that consists of workers
	workerPool *pool

	//queue to hold waiting tasks
	taskQueue  *taskQueue

	//size of the pool
	poolSize   int
}

//returns PoolService object with the default pool size
func NewDefaultSizePool() *PoolService {

	poolService := &PoolService{
		taskQueue: NewTaskQueue(),
		poolSize:  DefaultPoolSize,
	}
	pool := newFixedSizePool(DefaultPoolSize, poolService)
	poolService.workerPool = pool

	return poolService
}

//returns PoolService object with the specified pool size
func NewFixedSizePool(nGoRoutines int) (*PoolService, error) {

	if nGoRoutines <= 0 {
		panic(ErrInvalidPoolSize)
	}
	poolService := &PoolService{
		taskQueue: NewTaskQueue(),
		poolSize:  nGoRoutines,
	}
	pool := newFixedSizePool(nGoRoutines, poolService)
	poolService.workerPool = pool

	return poolService, nil
}

//submits a new task and assigns it to the pool
func (rs *PoolService) Submit(fun func()) error {
	if fun == nil {
		return nil
	}

	if rs.workerPool.status == CLOSED {
		return ErrInvalidPoolState
	}

	newTask := Task{executable: fun}

	//	if worker is available, immediately assigning the task
	if rs.workerPool.isWorkerAvailable() {
		fmt.Println("Assigning!")
		rs.workerPool.assignTask(newTask)
	} else {
		fmt.Println("Queuing!")
		rs.taskQueue.EnqueueTask(newTask)
	}
	return nil
}

//notifies the poolService that one of the worker from
//the pool has completed its task and a new task can be
//assigned to this worker if waiting in the queue.
func (rs *PoolService) notify() {
	if rs.taskQueue.IsNotEmpty() {
		task, err := rs.taskQueue.DequeueTask()
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error Dequeueing Task!")
			return
		}
		fmt.Println("Pulling task from Queue!")
		rs.workerPool.assignTask(task)
	}
}

//returns active workers
func (rs *PoolService) ActiveWorkers() int {
	return rs.workerPool.activeWorkers
}

//returns pool size
func (rs *PoolService) PoolSize() int {
	return rs.workerPool.poolSize
}

//returns available workers out of total workers
func (rs *PoolService) AvailableWorkers() int {
	return rs.workerPool.poolSize - rs.workerPool.activeWorkers
}

//closes the pool
func (rs *PoolService) Close() {
	rs.workerPool.close()
}
