package gohive

import "log"

//Represents a worker which is responsible to execute a task,
//handle panics if any and after completing the task, notifies
//back to the pool to pull new task from the TaskQueue
type worker struct {
	//channel that receives a task
	taskChan chan func()

	//reference back to the pool who owns this worker
	pool *pool
}

//This method executes the task
func (w *worker) run() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered Error: %v", r)
				w.pool.done(w)
			}
		}()
		fun := <-w.taskChan
		fun()
		w.pool.done(w)
	}()
}
