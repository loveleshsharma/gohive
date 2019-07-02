package GoHive

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDefaultSizePool(t *testing.T) {

	testPoolService := NewDefaultSizePool()

	assert.NotNil(t, testPoolService, "pool service should not be nil!")
	assert.Equal(t, 10, testPoolService.workerPool.poolCapacity, "default pool size should be 10!")

}

func TestNewFixedSizePool(t *testing.T) {

	testPoolService := NewFixedSizePool(15)

	assert.NotNil(t, testPoolService, "pool service should not be nil!")
	assert.Equal(t, 15, testPoolService.poolSize, "pool size should be 15!")

}

func TestNewFixedSizePoolWhenPoolSizeIsNegative(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Invalid pool size: pool size must be a positive number!", r.(error).Error())
		}
	}()
	testPoolService := NewFixedSizePool(-1)
	assert.Nil(t, testPoolService)

}

func TestPoolService_Submit_WhenFunIsNil(t *testing.T) {

	testPoolService := NewFixedSizePool(10)

	err := testPoolService.Submit(nil)

	assert.Equal(t, "Cannot submit Nil function()!", err.Error())

}

func TestPoolService_Submit_WhenPoolIsClosed(t *testing.T) {

	testPoolService := NewFixedSizePool(10)
	testPoolService.Close()

	err := testPoolService.Submit(func() {
		fmt.Println("Test Function")
	})

	assert.Equal(t, "Pool is closed: cannot assign task to a closed pool!", err.Error())

}

func TestPoolService_Submit_WhenPoolIsFull(t *testing.T) {

	testPoolService := NewFixedSizePool(1)

	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
	})
	testPoolService.Submit(func() {
		fmt.Println("second function")
	})

	assert.Equal(t, 1, testPoolService.taskQueue.totalTasks)

}

func TestPoolService_Submit(t *testing.T) {

	testPoolService := NewDefaultSizePool()

	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Test Function")
	})

	assert.Equal(t, 1, testPoolService.ActiveWorkers())

}

func TestPoolService_notify(t *testing.T) {

	testPoolService := NewFixedSizePool(1)

	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("First Function")
	})
	testPoolService.Submit(func() {
		fmt.Println("Second Function")
	})

	assert.Equal(t, 1, testPoolService.taskQueue.totalTasks)

	time.Sleep(2 * time.Second)

	assert.Equal(t, 0, testPoolService.taskQueue.totalTasks)

}

func TestPoolService_ActiveWorkers(t *testing.T) {

	testPoolService := NewFixedSizePool(5)

	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("First Function")
	})
	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Second Function")
	})

	assert.Equal(t, 2, testPoolService.ActiveWorkers())

}

func TestPoolService_AvailableWorkers(t *testing.T) {

	testPoolService := NewFixedSizePool(5)

	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("First Function")
	})
	testPoolService.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Second Function")
	})

	assert.Equal(t, 3, testPoolService.AvailableWorkers())

}

func TestPoolService_PoolSize(t *testing.T) {

	testPoolService := NewFixedSizePool(5)

	assert.Equal(t, 5, testPoolService.PoolSize())

}
