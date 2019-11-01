package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 使用有缓冲通道和固定数目的goroutine来处理一堆工作

const (
	numberGoroutines = 4	// 要使用的goroutine数量
	taskLoad = 10 	// 任务数量
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().Unix())
}


func main() {

	// 创建有缓冲通道来管理任务
	tasks := make(chan string, taskLoad)

	// 启动goroutine处理任务
	wg.Add(numberGoroutines)
	for gr:=1; gr<=numberGoroutines; gr++ {
		go worker(tasks, gr)
	}

	// 增加一组要完成的工作
	for post:=1; post<=taskLoad; post++ {
		tasks <- fmt.Sprintf("Task: %d", post)
	}

	// 所有任务完成后关闭通道，以便所有goroutine退出
	close(tasks)		// 注意！！！ 关闭通道后，仍可以从通道接收（这允许通道关闭后仍能取出所有缓冲数据），但不可发送。		另外，不能关闭一个仅可接收的单向通道

	// 等待所有工作完成
	wg.Wait()
}

// worker作为goroutine启动来处理从有缓冲通道传入的工作
func worker(tasks chan string, workerNo int) {

	defer wg.Done()

	for {

		// 等待分配工作
		task, ok := <- tasks
		if !ok {
			// 意味着通道已经空了并且被关闭
			fmt.Printf("Worker %d : shutting down\n", workerNo)
			return
		}

		// 显示已经工作，并模拟工作
		fmt.Printf("Worker %d : Started %s\n", workerNo, task)
		time.Sleep(time.Duration(rand.Int63n(100))*time.Millisecond)

		// 显示已经完成工作
		fmt.Printf("Worker %d : Completed %s\n", workerNo, task)

	}

}