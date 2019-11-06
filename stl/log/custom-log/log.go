package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace *log.Logger	// 记录所有日志
	Info *log.Logger	// 重要的信息
	Warning *log.Logger	// 需要注意的信息
	Error *log.Logger	// 非常严重的问题
)

func init() {
	// 打开errors.txt文档，为写入错误做准备
	file, err := os.OpenFile("errors.txt", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	// 定制各个logger
	Trace = log.New(ioutil.Discard, "Trace: ", log.Ldate | log.Ltime | log.Lshortfile)	// ioutil.Discard也是一个io.Writer但会将所有写入的内容丢弃。
	Info = log.New(os.Stdout, "Info: ", log.Ldate | log.Ltime | log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning: ", log.Ldate | log.Ltime | log.Lshortfile)
	Trace = log.New(io.MultiWriter(file, os.Stderr), "Error: ", log.Ldate | log.Ltime | log.Lshortfile)	// io.MultiWriter将多个writer绑定为一个，实现批量写

}

func main() {
	Trace.Println("I have something standard to say")
	Info.Println("Special Information")
	Warning.Println("There is something you need to know about")
	Error.Println("Something has failed")
}
