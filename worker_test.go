package gohive

import (
	"fmt"
	"testing"
)

func Test_worker_run(t *testing.T) {

	testPoolService := NewFixedSizePool(1)

	testPoolService.Submit(func() {
		a := 10
		b := 0
		//dividing by zero to raise a panic
		fmt.Println("Division: ", a/b)
	})

}
