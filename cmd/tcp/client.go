package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2020")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	//开启一个等待服务关闭的协程
	go func() {
		io.Copy(os.Stdout, conn) //从连接读到标准输出
		log.Println("done")
		done <- struct{}{} //给主协程发送信号
	}()
	mustCopy(conn, os.Stdin) //通过标准输入写入连接
	conn.Close()
	<-done

}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
