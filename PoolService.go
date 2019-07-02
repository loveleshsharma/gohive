package GoHive

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

//default pool size for number of workers in the pool
const DefaultPoolSize = 10

//error for invalid pool size
var ErrInvalidPoolSize = errors.New("Invalid pool size: pool size must be a positive number!")

//error for invalid pool state
var	ErrInvalidPoolState = errors.New("Pool is closed: cannot assign task to a closed pool!")

//error of nil function submitted
var ErrNilFunction = errors.New("Cannot submit Nil function()!")


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
	pool := newPool(DefaultPoolSize, poolService)
	poolService.workerPool = pool

	return poolService
}

//returns PoolService object with the specified pool size
func NewFixedSizePool(nGoRoutines int) *PoolService {

	if nGoRoutines <= 0 {
		panic(ErrInvalidPoolSize)
	}
	poolService := &PoolService{
		taskQueue: NewTaskQueue(),
		poolSize:  nGoRoutines,
	}
	pool := newPool(nGoRoutines, poolService)
	poolService.workerPool = pool

	return poolService
}

//submits a new task and assigns it to the pool
func (rs *PoolService) Submit(fun func()) error {
	if fun == nil {
		return ErrNilFunction
	}

	if rs.workerPool.status == CLOSED {
		return ErrInvalidPoolState
	}

	newTask := Task{executable: fun}

	//	if worker is available, immediately assigning the task
	if rs.workerPool.isWorkerAvailable() {
		rs.workerPool.assignTask(newTask)
	} else {
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
		rs.workerPool.assignTask(task)
	}
}

//returns active workers
func (rs *PoolService) ActiveWorkers() int {
	return rs.workerPool.getActiveWorkers()
}

//returns pool size
func (rs *PoolService) PoolSize() int {
	return rs.workerPool.getPoolSize()
}

//returns available workers out of total workers
func (rs *PoolService) AvailableWorkers() int {
	return rs.workerPool.getPoolSize() - rs.workerPool.getActiveWorkers()
}

//closes the pool
func (rs *PoolService) Close() {
	rs.workerPool.close()
}
