package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:2020")
	if err != nil {
		panic(err)
	}
	//广播
	go broadCaster()

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Println("错误：%v", err)
			continue
		}
		go handleConn(conn)

	}
}

// User 用户信息
type User struct {
	ID            string
	Address       string
	EnterAt       time.Time
	MessageChanel chan string
}

var (
	messageChanel  = make(chan string, 8) //全局通信通道
	enteringChanel = make(chan *User)     //进入的用户
	leavingChanel  = make(chan *User)     //退出的用户

)

// broadCast 记录用户信息，广播信息
// 1.新用户进入，2.普通用户消息，3.用户退出
func broadCaster() {
	var users = make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChanel:
			//新用户进入
			users[user] = struct{}{}
		case user := <-leavingChanel:
			//用户退出，删除用户
			delete(users, user)
			//关闭通道，防止goroutine溢出
			close(user.MessageChanel)
		case message := <-messageChanel:
			for user := range users {
				//给其他用户发送消息
				user.MessageChanel <- message
			}
		}
	}

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
	user.MessageChanel <- "welcome user: " + user.String()
	messageChanel <- "user:`" + user.ID + "` entering"
	//4.将新进入的用户加入全部用户列表中，避免加锁
	enteringChanel <- user
	//5.循环读取用户发送的消息
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChanel <- user.ID + ":" + input.Text()
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误:%v", err)
	}
	//6.用户退出
	leavingChanel <- user
	messageChanel <- "user:" + user.ID + "has left"
}

// sendMessage 给其他用户发送消息
func sendMessage(conn net.Conn, messageChanel <-chan string) {
	for message := range messageChanel {
		fmt.Fprintln(conn, message)
	}
}

// GenUserId 生成用户id
func GenUserId() string {
	return uuid.New().String()
}
func (user *User) String() string {
	return " UID:" + user.ID + " address: " + user.Address + " enterAt " + user.EnterAt.String()
}
