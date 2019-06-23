package GoHive

import (
	"fmt"
	"os"
)

type worker struct {
	taskChan chan func()
	pool     *pool //reference back to the pool who owns this worker
}

func (w *worker) run() {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				fmt.Fprintf(os.Stderr, "Recovered Error: %v", err)
			}
			w.pool.done(w)
		}
	}()
	fun := <-w.taskChan
	fun()
	w.pool.done(w)
}
