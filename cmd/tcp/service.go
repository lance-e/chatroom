package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// User 用户信息
type User struct {
	ID            int
	Address       string
	EnterAt       time.Time
	MessageChanel chan string
}

var messageChanel chan string //全局通信通道
var enteringChanel chan *User //进入的用户
var leavingChanel chan *User  //退出的用户
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:2020")
	if err != nil {
		panic(err)
	}
	go broadCast()

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("错误：%v", err)
			continue
		}
		go handleConn(conn)

	}
}

// broadCast 记录用户信息，广播信息
// 新用户进入，普通用户，用户退出
func broadCast() {

}
func handleConn(conn net.Conn) {
	defer conn.Close()
	//1.初始化一个用户实例
	var user = &User{
		ID:            GenUserId(),
		Address:       conn.RemoteAddr().String(),
		EnterAt:       time.Now(),
		MessageChanel: make(chan string, 8),
	}
	//2.启动一个给用户发送消息的协程，
	go sendMessage(conn, user.MessageChanel)
	//3.给新进入的用户发送信息，以及给其他用户通知
	user.MessageChanel <- "welcome " + user.String()
	messageChanel <- "user:" + strconv.Itoa(user.ID) + "entering"
	//4.将新进入的用户加入全部用户列表中，避免加锁
	enteringChanel <- user
	//5.循环读取用户发送的消息
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChanel <- strconv.Itoa(user.ID) + ":" + input.Text()
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误:%v", err)
	}
	//6.用户退出
	leavingChanel <- user
	messageChanel <- "user:" + strconv.Itoa(user.ID) + "has left"
}

// sendMessage 给其他用户发送消息
func sendMessage(conn net.Conn, messageChanel <-chan string) {
	for message := range messageChanel {
		fmt.Sprintf("message:%v", message)
	}
}
