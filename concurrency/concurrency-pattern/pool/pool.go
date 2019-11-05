// 使用带缓冲通道实现资源池，管理可以在任意多个数量的goroutine之间共享或独立使用的资源。
// 这种模式在需要共享一组静态资源的情况（比如共享数据库连接或者内存缓冲区）时非常有用。
// 如果goroutine需要使用资源池内资源，需要从池里申请，使用完归还到资源池
// 现在go SDK已经自带sync.Pool这样一个资源池的实现。

// 资源，比如说和数据库建立的连接，如果每一次使用数据库，都重新建一次连接，那么开销比较大，可以通过资源池的方法来让其尽可能复用。

package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

// Pool 管理一组可以安全的在多个goroutine间共享的资源。被管理的资源必须实现io.Closer接口
type Pool struct {
	m sync.Mutex	// 互斥锁
	resources chan io.Closer		// 管理的资源
	factory func() (io.Closer, error)	// 工厂方法，由外部实现并传入，用来创建由池管理的资源
	closed bool		// 标志pool是否被关闭
}

// ErrPoolClosed 表示请求（acquire）了一个已经关闭的Pool
var ErrPoolClosed = errors.New("pool has been closed")

// New 创建一个用于管理资源的池
// 这个池需要一个可以分配新资源的函数，并规定池的大小
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size value too small")
	}

	return &Pool{
		factory:   fn,
		resources:make(chan io.Closer, size),
	}, nil
}

// Acquire 从池中获取一个资源
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	// 检查是否有空闲资源
	case r, ok := <-p.resources:
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	// 没有可用资源的话，就提供一个新资源
	default:
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}

// Release 将一个使用过的资源放回池中
func (p *Pool) Release(r io.Closer) {
	// 保证本操作和Close操作的安全。使得在不同goroutine间close()和release不能同时进行。
	p.m.Lock()
	defer p.m.Unlock()

	// 被保护的操作

	// 若池已被关闭，则销毁这个资源
	if p.closed {
		r.Close()
		return
	}

	select {
	// 试图放回资源
	case p.resources<- r:
		log.Println("Release:", "In Queue")

	// 若队列已满，就放不成功，那么关闭这个资源
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

// Close 让资源池停止工作，并关闭所有现有资源
func (p *Pool) Close() {
	// 保证本操作与Release操作的安全
	p.m.Lock()
	defer p.m.Unlock()

	// 如果pool已经被关闭，什么也不做
	if p.closed {
		return
	}

	// 将池关闭
	p.closed = true

	// 清空通道里的资源之前，将通道关闭。 如果不这么做的话，会发生死锁
	close(p.resources)

	// 关闭资源
	for r := range p.resources {
		r.Close()
	}
}