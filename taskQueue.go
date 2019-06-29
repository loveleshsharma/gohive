package GoHive

import (
	"github.com/pkg/errors"
)

//default size of the taskQueue
const DefaultQueueSize = 10

//represents a queue that holds tasks which
//are in waiting for workers from the pool
type taskQueue struct {

	//queue that holds tasks
	que        []Task

	//number of tasks that currently resides in the queue
	totalTasks int
}

//returns new taskQueue with the default capacity
func NewTaskQueue() *taskQueue {
	wtQue := taskQueue{que: make([]Task, 0, DefaultQueueSize), totalTasks: 0}
	return &wtQue
}

//puts a new task in the taskQueue
func (wq *taskQueue) EnqueueTask(task Task) {
	wq.que = append(wq.que, task)
	wq.totalTasks++
}

//returns a task and removes it from the taskQueue
func (wq *taskQueue) DequeueTask() (Task, error) {
	if wq.totalTasks > 0 {
		task := wq.que[0]
		wq.que = append(wq.que[:0], wq.que[1:]...)
		wq.totalTasks--
		return task, nil
	} else {
		return Task{}, errors.New("Queue is Empty")
	}
}

//This function returns whether the taskQueue is empty or not
func (wq *taskQueue) IsNotEmpty() bool {
	return wq.totalTasks > 0
}