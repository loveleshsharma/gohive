# GoHive

<img src="https://github.com/loveleshsharma/GoHive/blob/master/beehive.jpg"/>

Package GoHive implements a simple and easy to use goroutine pool for Go

## Features

- Pool can be created with a specific size as per the requirement
- Offers efficient performance by implementing ```sync.Pool```, which maintains pool of workers in which workers gets recycled automatically when not in use  
- Implements a <B>Task Queue</B> which can hold surplus tasks in waiting, if submitted more than the pool capacity
- Implements PoolService type, which acts as an easy to use API with simple methods to interact with GoHive
- Gracefully handles panics and prevent the application from getting crashed or into deadlocks
- Provides functions like: AvailableWorkers(), ActiveWorkers() and Close() etc.

## Installation
```go
go get -u github.com/loveleshsharma/GoHive
```

## Usage

- Create an instance of PoolService type first

```go
hive := GoHive.NewFixedSizePool(5)
```

- Invoke the Submit() function and pass the task to execute

```
hive.Submit(someTask())
```
Submit function accepts a function as an argument, which it passes to the pool if a worker is available, otherwise enqueues it in a waiting queue

- To close the pool we can invoke the Close() function

```
hive.Close()
```
Once the pool is closed, we cannot assign any task to it

## Example

Let's get into a full program where we can see how to use the GoHive package in order to execute many goroutines simultaneously

```go
package main

import (
	"github.com/loveleshsharma/GoHive"
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	hivePool := GoHive.NewFixedSizePool(5)

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