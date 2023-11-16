package gohive

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testRunnable struct{}

func (t testRunnable) Run() {
	for {
		time.Sleep(time.Millisecond * 100)
	}
}

var runnableObject = testRunnable{}

func TestPool_CloseShouldReturnErrorIfAlreadyClosed(t *testing.T) {
	testPool := NewFixedPool(5)

	_ = testPool.Close()

	actualError := testPool.Close()

	assert.NotNil(t, actualError, "error should be not nil")

}

func TestPool_CloseShouldReturnNilIfPoolIsOpen(t *testing.T) {
	testPool := NewFixedPool(5)

	actualError := testPool.Close()

	assert.Nil(t, actualError, "error should be nil")
}

func TestPool_IsPoolClosedShouldReturnTrueOrFalseIfThePoolIsClosed(t *testing.T) {
	testPool := NewFixedPool(5)

	closed := testPool.IsPoolClosed()
	assert.Falsef(t, closed, "pool should not be closed")

	_ = testPool.Close()
	closed = testPool.IsPoolClosed()
	assert.True(t, closed, "pool should be closed")
}

func TestPool_SubmitShouldReturnErrorIfRunnableIsPassedAsNil(t *testing.T) {
	testPool := NewFixedPool(5)

	actualError := testPool.Submit(nil)

	assert.NotNil(t, actualError, "Submit should return error if runnable is nil")
}

func TestPool_SubmitShouldReturnErrorIfPoolIsClosed(t *testing.T) {
	testPool := NewFixedPool(5)

	_ = testPool.Close()
	actualError := testPool.Submit(runnableObject)

	assert.NotNil(t, actualError, "Submit should return error if runnable is nil")
}
