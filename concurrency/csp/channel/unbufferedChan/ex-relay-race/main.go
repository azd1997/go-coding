package main

import (
	"fmt"
	"sync"
	"time"
)

// 无缓冲通道模拟四个goroutine间的接力比赛


var wg sync.WaitGroup

func main() {

	// 模拟接力棒
	baton := make(chan int)

	// 为最后一个跑步者计数加1，等待它跑完
	wg.Add(1)

	// 第一位跑步者准备持有接力棒
	go runner(baton)

	// 开始比赛，接力棒发出
	baton <- 1

	// 等待比赛结果
	wg.Wait()

}

// runner模拟接力赛中的跑步者
func runner(baton chan int) {

	var newRunner int

	// 等待接力棒
	_runner := <-baton

	// 开始跑步
	fmt.Printf("Runner %d Running With Baton\n", _runner)

	// 创建下一位跑步者
	if _runner != 4 {	// 如果当前runner不是第四个，他就需要把接力棒传给下一位
		newRunner = _runner + 1
		fmt.Printf("Runner %d To The Line\n", newRunner)		// 下一位跑步者到起跑线准备接棒
		go runner(baton)
	}

	// 模拟跑步过程
	time.Sleep(100*time.Millisecond)

	// 比赛结束了吗？
	if _runner == 4 {
		fmt.Printf("Runner %d Finished, Race Over\n", _runner)
		wg.Done()
		return
	}

	// 没结束的话将接力棒交给下一个人
	fmt.Printf("Runner %d Exchange With Runner %d\n", _runner, newRunner)
	baton <- newRunner
}

