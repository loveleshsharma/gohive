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
- Accepts tasks which implements Runner interface
- Uses channels to accepts tasks and gets them executed via workers
- Uses synchronization among workers to avoid race conditions

## Installation
Use ```go get``` to install and update:
```go
$ go get -u github.com/loveleshsharma/gohive
```

## Usage

- Create an instance of Pool type first

```go
hive := gohive.NewFixedPool(5)
```

- Invoke the Submit() function and pass the task to execute

```go
hive.Submit(object Runner)
```
Submit function accepts a Runner object as an argument, which it passes to the pool if a worker is available, otherwise it will wait for the worker to be available

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
   "fmt"
   "github.com/loveleshsharma/gohive"
   "sync"
)

func main() {
   var wg sync.WaitGroup
   pool := gohive.NewFixedPool(5)

   for i := 1; i <= 20; i++ {
      if err := pool.Submit(NewMyStruct(i, &wg)); err != nil {
         fmt.Println("error: ", err)
         break
      }
   }

   wg.Wait()
}

type MyStruct struct {
   num int
   wg  *sync.WaitGroup
}

func NewMyStruct(num int, wg *sync.WaitGroup) MyStruct {
   myStruct := MyStruct{
      num: num,
      wg:  wg,
   }
   wg.Add(1)
   return myStruct
}

func (s MyStruct) Run() {
   defer s.wg.Done()
   val := s.num
   fact := s.num
   for i := s.num - 1; i > 0; i-- {
      fact *= i
   }

   fmt.Printf("Factorial of %d: %d\n", val, fact)
}


```
<B>Important : </B> Always keep sync.WaitGroup in your struct and put ```defer wg.Done()``` as the first statement of your Run() function. It will wait for your task to complete.

###
TODO
1. Maintain a waiting queue to stop blocking submit method when all goroutines are busy.
2. Submitting priority tasks which takes priority over other tasks.
3. Handling panics inside goroutines to prevent them from crashing.
4. Implement dynamic pool which will scale the number of goroutines as per requirement and scales down when they are idle.
5. Submitting multiple tasks together.