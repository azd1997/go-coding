package main

import (
	"fmt"
	"runtime"
	"sync"
)

// mutex: mutual exclusion
// 用于定义一段需要同步访问的代码临界区资源的同步访问
// 临界区和原子操作的概念不同，原子操作底层绑定一个硬件级别的操作，不可中断。而临界区则是用户自己前后上锁，中间这块叫临界区

// 全局计数器

var (
	counter int	 // 每个goroutine都要增加其值
	wg sync.WaitGroup
	mutex sync.Mutex	// 互斥锁，用来定义一段临界区
)

func main() {

	wg.Add(2)

	go incCounter(1)
	go incCounter(2)

	// 等待goroutine结束
	wg.Wait()
	fmt.Printf("Final Counter: %d \n", counter)		// 最终为4
}

// 使用互斥锁来同步并保证counter的安全访问（增加counter的值）
func incCounter(id int) {
	defer wg.Done()

	for count:=0; count<2; count++ {

		// 上锁
		mutex.Lock()

		// 锁住的时间内为临界区。{}不是必需，但可以让代码清晰
		{
			// 捕获counter值
			value := counter

			// 当前goroutine从线程退出，并放回队列. (就是放弃当着这次调度执行的权利，下次再调度)
			// 这是故意让这个工作执行一半后就切到另一个工作，来看看go scheduler的作用
			runtime.Gosched()

			// 增加本地value值
			value++

			// 将该值保存回counter
			counter = value
		}


		// 解锁
		mutex.Unlock()
	}
}

