package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// unbufferedChan 指在接收前没有能力保存任何值的通道。要求发送方goroutine和接收方goroutine同时准备好，才能完成发送和接收操作。
// 如果一方没有准备好，先执行的那一方会阻塞直到对方准备好。
// 这是一种同步通信。任意一个操作无法离开另一个操作


// 该代码使用两个goroutine模拟网球比赛，两个选手总是处于等待接球和接到球将球击回两个动作之一，无缓冲通道模拟球的传递


var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())	// 使用当前时间的UnixNano戳作为随机数的生成源
}

func main() {

	// 创建无缓冲通道
	court := make(chan int)

	// 计数器+2，等待两个goroutine
	wg.Add(2)

	// 启动两个goroutine(选手)
	go player("Alice", court)
	go player("Bob", court)

	// 发球
	court <- 1

	// 等待两个goroutine结束（游戏结束）
	wg.Wait()
}

// 模拟一个网球选手在打球
func player(name string, court chan int) {

	defer wg.Done()

	for {
		// 等待球被击过来（chan接收到数据）
		ball, ok := <- court
		if !ok {
			// 如果通道关闭，接收不到数据，说明对方没能把球击回来，那就说明赢了。那么当前函数就返回，比赛结束
			fmt.Printf("Player %s Won\n", name)
			return
		}

		// 选随机数，用这个随机数判断是否丢球（也就是比赛时没接住对方的球），就得将通道关闭
		n := rand.Intn(100)	// 0~99随机整数
		if n%13 == 0 {
			fmt.Printf("Player %s Missed\n", name)

			// 关闭通道然后退出
			close(court)
			return
		}

		// 显示击球数，并将击球数加1
		fmt.Printf("Player %s Hits %d \n", name, ball)
		ball++

		// 将球打向对手
		court <- ball	// 这个时候两个goroutine都会锁住，直到数据传输完成

	}
}