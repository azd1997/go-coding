/*
work包展示如何使用无缓冲通道实现一个goroutine池，这些goroutine执行并控制一组工作，让其并发执行。
这种情形下，使用无缓冲通道比虽已制定一个缓冲区大小的有缓冲通道好，因为这个情况下既不需要一个工作队列，也不需要一组goroutine配合执行。

*/
package work


