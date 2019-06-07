package GoHive

type Pool struct {
	capacity       int
	runningWorkers int
	workers        []Worker
	routineService *RoutineService  //reference back to the routine service
}

func NewFixedSizePool(newSize int, routineService *RoutineService) Pool {
	newPool := Pool{capacity: newSize, workers: make([]Worker, newSize), routineService: routineService}
	return newPool
}

func (p *Pool) assignTask(task Task) {
	for i := 0; i < p.capacity; i++ {
		if p.workers[i].isOccupied == false {
			p.workers[i] = Worker{taskChan: make(chan func()),pool:p,isOccupied:true}
			go p.workers[i].run()
			p.workers[i].taskChan <- task.getTask()
			p.runningWorkers++
			break
		}
	}
}

func (p *Pool) Done(w *Worker) {
	p.runningWorkers--
	w.isOccupied = false
	p.routineService.notify() 	//	notify the RoutineService to pull more tasks from the queue if waiting
}

func (p *Pool) isWorkerAvailable() bool {
	return p.capacity > p.runningWorkers
}

func (p *Pool) getRunning() int {
	return p.runningWorkers
}

func (p *Pool) getCapacity() int {
	return p.capacity
}
