package GoHive

type Worker struct {
	isOccupied bool
	taskChan   chan func()
}

func (w *Worker) run() {
	fun := <- w.taskChan
	fun()
}
