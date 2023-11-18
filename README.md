<p align="center"> 
    <img width="150" height="150" src="internal/static/GoHiveIcon.png" alt="">
    <h1 align="center">gohive</h1>
    <p align="center">Package gohive implements a simple and easy to use goroutine pool for Go<p>
    <p align="center"><a href="https://travis-ci.org/loveleshsharma/gohive"><img src="https://travis-ci.org/loveleshsharma/gohive.svg?branch=master" /></a>
    <a href="https://goreportcard.com/report/github.com/loveleshsharma/gohive"><img src="https://goreportcard.com/badge/github.com/loveleshsharma/gohive" /></a>
    <a href="https://codecov.io/gh/loveleshsharma/gohive"><img src="https://codecov.io/gh/loveleshsharma/gohive/branch/master/graph/badge.svg" /></a>
    <a href="https://godoc.org/github.com/loveleshsharma/gohive"><img src="https://godoc.org/github.com/loveleshsharma/gohive?status.svg" /></a>
    <a href="https://github.com/avelino/awesome-go#goroutines"><img src="https://awesome.re/mentioned-badge.svg" /></a>
    </p>
</p>

## Features

- Pool can be created with a specific size as per the requirement
- Offers efficient performance by implementing ```sync.Pool```, which maintains pool of workers in which workers gets recycled automatically when not in use  
- Implements a <B>Task Queue</B> which can hold surplus tasks in waiting, if submitted more than the pool capacity
- Implements PoolService type, which acts as an easy to use API with simple methods to interact with gohive
- Gracefully handles panics and prevent the application from getting crashed or into deadlocks
- Provides functions like: AvailableWorkers(), ActiveWorkers() and Close() etc.

## Installation
Use ```go get``` to install and update:
```go
$ go get -u github.com/loveleshsharma/gohive
```

## Usage

- Create an instance of PoolService type first

```go
hive := gohive.NewFixedSizePool(5)
```

- Invoke the Submit() function and pass the task to execute

```go
hive.Submit(someTask())
```
Submit function accepts a function as an argument, which it passes to the pool if a worker is available, otherwise enqueues it in a waiting queue

- To close the pool we can invoke the Close() function

```go
hive.Close()
```
Once the pool is closed, we cannot assign any task to it

## Example

Let's get into a full program where we can see how to use the gohive package in order to execute many goroutines simultaneously

```go
package main

import (
	"github.com/loveleshsharma/gohive"
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	hivePool := gohive.NewFixedSizePool(5)

	//wrap your executable function into another function with wg.Done()
	executableTask := func() {
		defer wg.Done()
		factorial(5)
	}

	wg.Add(1)
	hivePool.Submit(executableTask)

	wg.Wait()
}

func factorial(val int) {
	var fact = val
	var res = 1

	for i := fact; i > 0; i-- {
		res = res * i
	}

	fmt.Printf("Factorial: %v", res)
}

```
<B>Important : </B> Always put ```defer wg.Done()``` as the first statement of your wrapper function. It will prevent your program from deadlocks in case of panics

Workers implements a notifying mechanism, due to which they can notify to the pool that their task is completed and they are available to execute more tasks if in waiting queue

###
TODO
1. Maintain a waiting queue to stop blocking submit method when all goroutines are busy.
2. Submitting priority tasks which takes priority over other tasks.
3. Handling panics inside goroutines to prevent them from crashing.
4. Implement dynamic pool which will scale the number of goroutines as per requirement and scales down when they are idle.
5. Submitting multiple tasks together.