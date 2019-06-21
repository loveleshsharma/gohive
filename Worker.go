package GoHive

type Worker struct {
	taskChan   chan func()
	pool 	   *Pool	//reference back to the pool who owns this worker
}

func (w *Worker) run() {
	fun := <- w.taskChan
	fun()
	w.pool.Done(w)
}
