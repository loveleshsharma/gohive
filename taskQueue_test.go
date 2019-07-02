package GoHive

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_taskQueue_EnqueueTask(t *testing.T) {

	testQueue := NewTaskQueue()

	testQueue.EnqueueTask(NewTask(func() {
		fmt.Println("sample func")
	}))

	assert.Equal(t, 1, testQueue.totalTasks, "task count in the queue shoule be 1!")

	testQueue.EnqueueTask(NewTask(func() {
		fmt.Println("sample func")
	}))

	assert.Equal(t, 2, testQueue.totalTasks, "task count in the queue shoule be 2!")

}

func Test_taskQueue_DequeueTask(t *testing.T) {

	testQueue := NewTaskQueue()

	testQueue.EnqueueTask(NewTask(func() {
		fmt.Println("sample func")
	}))

	assert.Equal(t, 1, testQueue.totalTasks, "task count in the queue should be 1!")

	testQueue.DequeueTask()

	assert.Equal(t, 0, testQueue.totalTasks, "task count in the queue should be 0!")
}

func Test_taskQueue_DequeueTaskWhenQueueIsEmpty(t *testing.T) {

	testQueue := NewTaskQueue()

	_, err := testQueue.DequeueTask()

	assert.Equal(t,"Queue is Empty",err.Error(),"Error message should be correct!")

}

func Test_taskQueue_IsNotEmpty(t *testing.T) {

	testQueue := NewTaskQueue()

	testQueue.EnqueueTask(NewTask(func() {
		fmt.Println("sample func")
	}))

	assert.Equal(t, true, testQueue.IsNotEmpty(), "taskQueue should not be empty!")

	testQueue.DequeueTask()

	assert.Equal(t, false, testQueue.IsNotEmpty(), "taskQueue should be empty!")
}
