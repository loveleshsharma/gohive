package GoHive

type Pool struct {
	capacity int
	running  int
	workers  []Worker
}

func NewFixedSizePool(newSize int) Pool {
	newPool := Pool{capacity: newSize, workers: make([]Worker, newSize)}
	return newPool
}

func (p *Pool) assignTask(task Task) {
	for i := 0; i < p.capacity; i++ {
		if p.workers[i].isOccupied == false {
			p.workers[i] = Worker{taskChan:make(chan func())}
			go p.workers[i].run()
			p.workers[i].taskChan <- task.getTask()
			break
		}
	}
}

func (p *Pool) isWorkerAvailable() bool {
	return p.capacity > p.running
}
