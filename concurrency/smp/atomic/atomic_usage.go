package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	shutdown int64	// 通知正在运行的goroutine停止的标志
	wg sync.WaitGroup	// 用来等待所有工作结束
)

func main() {

	wg.Add(2)		// 表示要等待2个goroutine

	// 创建两个goroutine
	go doWork("A")
	go doWork("B")

	// 给定goroutine执行时间，在这睡一秒，然后就通过标志通知两个goroutine结束
	time.Sleep(1*time.Second)
	fmt.Println("Shutdown Now")
	atomic.StoreInt64(&shutdown, 1)		// atomic.StoreInt64用于原子写，常用的还有atomic.LoadInt64(读)和atomic.AddInt64(加)

	wg.Wait()	// 等待两goroutine结束，然后就退出程序
}

// 模拟执行工作的goroutine。如果检测到shutdown标志置1那么就提前结束，释放goroutine
func doWork(name string) {
	// 在函数退出前调用wg.Done告诉main goroutine 工作完成
	defer wg.Done()

	// 模拟循环工作
	for {
		// 模拟工作一次
		fmt.Printf("Doing %s work\n", name)
		time.Sleep(250*time.Millisecond)

		// 检测是否要关闭
		if atomic.LoadInt64(&shutdown) == 1 {
			fmt.Printf("Shutting %s down\n", name)
			break	// 退出工作程序，go调度器会将这个goroutine释放
		}
	}
}

