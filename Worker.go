package GoHive

type worker struct {
	taskChan   chan func()
	pool 	   *pool //reference back to the pool who owns this worker
}

func (w *worker) run() {
	defer func() {
		if r := recover(); r != nil {
			w.pool.done(w)
		}
	}()
	fun := <- w.taskChan
	fun()
	w.pool.done(w)
}
