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

	defer conn.Close()
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) //通过标准输出，从连接读取
		log.Println("done")
		done <- struct{}{}
	}()
	mustCopy(conn, os.Stdin) //通过标准输入写入连接
	<-done

}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
