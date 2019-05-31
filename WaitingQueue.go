package GoHive

import (
	"github.com/pkg/errors"
)

type WaitingQueue struct {
	que        []Task
	totalTasks int
}

func NewWaitingQueue() WaitingQueue {
	wtQue := WaitingQueue{que: make([]Task, 0, 10), totalTasks: 0}
	return wtQue
}

func (wq *WaitingQueue) EnqueueTask(task Task) {
	wq.que = append(wq.que, task)
	wq.totalTasks = len(wq.que)
}

func (wq *WaitingQueue) DequeueTask() (Task, error) {
	if len(wq.que) > 0 {
		task := wq.que[0]
		wq.que = append(wq.que[:0], wq.que[1:]...)
		wq.totalTasks = len(wq.que)
		return task, nil
	} else {
		return Task{}, errors.New("Queue is Empty")
	}
}
