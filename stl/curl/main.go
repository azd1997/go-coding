package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	// curl www.baidu.com ./file.txt

	// http get
	r, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// 创建一个文件 (os.Create(name))
	file, err := os.OpenFile(os.Args[2], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// 同时向文件和标准输出写入
	dest := io.MultiWriter(file, os.Stdout)
	io.Copy(dest, r.Body)
	// 关闭r.Body
	if err = r.Body.Close(); err != nil {
		log.Println(err)
	}
}
