package gohive

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

//DefaultPoolSize is a default size for number of workers in the pool
const DefaultPoolSize = 10

//ErrInvalidPoolSize indicates that the pool size is invalid
var ErrInvalidPoolSize = errors.New("Invalid pool size: pool size must be a positive number")

//ErrInvalidPoolState indicates that the invalid pool state
var ErrInvalidPoolState = errors.New("Pool is closed: cannot assign task to a closed pool")

//ErrNilFunction indicates that a nil function submitted
var ErrNilFunction = errors.New("Cannot submit Nil function()")

//PoolService acts as an orchestrator of the entire GoHive functionality
//It consists of a pool, that manages workers that run tasks and it
//consists of a TaskQueue that holds tasks waiting to acquire a worker
type PoolService struct {
	//pool that consists of workers
	workerPool *pool

	//queue to hold waiting tasks
	taskQueue *TaskQueue

	//size of the pool
	poolSize int
}

//NewDefaultSizePool returns PoolService object with the default pool size
func NewDefaultSizePool() *PoolService {

	poolService := &PoolService{
		taskQueue: NewTaskQueue(),
		poolSize:  DefaultPoolSize,
	}
	pool := newPool(DefaultPoolSize, poolService)
	poolService.workerPool = pool

	return poolService
}

//NewFixedSizePool returns PoolService object with the specified pool size
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

//Submit submits a new task and assigns it to the pool
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
			fmt.Fprintf(os.Stderr, "Error Dequeueing Task")
			return
		}
		rs.workerPool.assignTask(task)
	}
}

//ActiveWorkers returns active workers
func (rs *PoolService) ActiveWorkers() int {
	return rs.workerPool.getActiveWorkers()
}

//PoolSize returns pool size
func (rs *PoolService) PoolSize() int {
	return rs.workerPool.getPoolSize()
}

//AvailableWorkers returns available workers out of total workers
func (rs *PoolService) AvailableWorkers() int {
	return rs.workerPool.getPoolSize() - rs.workerPool.getActiveWorkers()
}

//Close closes the pool
func (rs *PoolService) Close() {
	rs.workerPool.close()
}
