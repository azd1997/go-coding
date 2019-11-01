// 用于展示如何使用通道来监视程序的执行时间，如果程序执行时间太长，还可以用runner包终止程序
// 当开发需要调度后台任务处理的程序时，这种病发模式很有用。
// 这个程序可能作为cron作业执行，或者是基于定时任务的云环境执行

// 实现了：
// 程序在超时时间内完成，正常终止；
// 超时未完成，“自杀”；
// 系统中断，程序清理状态并停止工作；
package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// Runner 在给定的超时时间内执行一组任务，并在操作系统发送中断信号时结束这些任务
type Runner struct {

	// interrupt 通道报告从操作系统发来的信号
	interrupt chan os.Signal

	// complete 通道报告处理任务已经完成
	complete chan error

	// timeout 报告处理任务已经超时
	timeout <- chan time.Time

	// tasks 持有一组以索引顺序依次执行的函数
	tasks []func(int)
}

// ErrTimeout 任务超时时返回
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt 会在接收到操作系统的事件时返回
var ErrInterrupt = errors.New("received interrupt")

// New 工厂方法，返回一个准备使用的Runner
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),		// 有缓冲，保证至少可以接收一个系统中断信号，并且发送这个事件时不会被阻塞
		complete:  make(chan error),
		timeout:   time.After(d),	// time.After(d)返回一个time.Time类型的通道，时间一到，就会从这个通道把时间发过来，起到定时通知的作用
		tasks:     nil,
	}
}

// Add 将一个任务附加到一个Runner上。任务是接收int类型的ID作为参数的函数
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start 执行所有任务，并且监视通道事件
func (r *Runner) Start() error {
	// 希望接收所有中断信号
	signal.Notify(r.interrupt, os.Interrupt)	// os/signal包，相当于os所有中断一旦产生都会通过通道传递给Runner

	// 开启goroutine执行任务
	go func() {
		r.complete <- r.run()	// r.run是真正拉起所有任务的方法
	}()

	select {
	// 任务完成时发送的信号
	case err := <-r.complete:
		return err
	// 任务超时的信号
	case <-r.timeout:
		return ErrTimeout
	}
}

// run 执行每一个已经注册的任务
func (r *Runner) run() error {
	for id, task := range r.tasks {
		// 监测系统中断信号， 如有，则返回中断错误
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		// 执行已注册的任务
		task(id)
	}
	return nil
}

// gotInterrupt 验证是否接收到操作系统的中断信号
func (r *Runner) gotInterrupt() bool {
	select {	// 如果select里边没有接收，也没有default就会阻塞，default分支使得interrupt通道由阻塞变成非阻塞
	case <- r.interrupt:
		// 停止接收后续信号，并返回true
		signal.Stop(r.interrupt)
		return true

	default:
		// 没有中断信号，继续正常运行
		return false
	}

}

