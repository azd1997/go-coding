/*
work包展示如何使用无缓冲通道实现一个goroutine池，这些goroutine执行并控制一组工作，让其并发执行。
这种情形下，使用无缓冲通道比虽已制定一个缓冲区大小的有缓冲通道好，因为这个情况下既不需要一个工作队列，也不需要一组goroutine配合执行。
无缓冲通道保证两个goroutine之间的数据交换，这种使用无缓冲通道的方法允许使用者知道什么时候goroutine池在执行工作。
而且，如果池中所有goroutine都忙，无法接受新的任务，能够通过通道及时通知调用者。
使用无缓冲通道不会有工作在队列中丢失或者卡住，所有工作都会被处理。
*/
package work

import "sync"

// Worker 工作者接口，符合这个接口才可以使用工作池
type Worker interface {
	Task()
}

// Pool 工作池，提供一个goroutine池，完成所有已提交的任务
type Pool struct {
	work chan Worker
	wg sync.WaitGroup
}

// New 新建工作池
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)

	for i:=0; i<maxGoroutines; i++ {	// 这做法就是先开这些协程，等待任务提交过来，提交过来的任务在池中会执行真正的任务过程
		go func() {
			for w := range p.work {		// 阻塞直到新任务从通道传过来。直到通道关闭，池关闭。
										// 这样的工作池可以不断的接收任务执行任务，重复利用已开的这些goroutine，避免了频繁的goroutine创建销毁过程。
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run 运行任务（提交任务到工作池）
func (p *Pool) Run(w Worker) {

	// 由于work是无缓冲通道，调用者（调用p.Run(w)）必须等待池中某个goroutine接收到这个值（这个工作）才会返回，否则一直阻塞。
	// 这样可以确保Run返回时，提交的工作已经开始执行。
	p.work<- w
}

// Shutdown 等待所有goroutine停止工作
func (p *Pool) Shutdown() {
	close(p.work)	// 关闭通道，不再接受新任务提交，但需要等待未完成任务继续执行完成
	p.wg.Wait()
}