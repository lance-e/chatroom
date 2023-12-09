package logic

import (
	"chatroom/global"
	"expvar"
	"fmt"
	"log"
)

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}
func calcMessageQueueLen() interface{} {
	fmt.Println("===len=:", len(BroadCaster.messageChannel))
	return len(BroadCaster.messageChannel)
}

// broadCast 广播器
type broadCaster struct {
	//所有聊天室用户
	users map[string]*User

	//所有goroutine统一管理，避免外部乱用

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message
	//判断用户昵称是否能够进入聊天室
	checkUserChannel      chan string //接收昵称，方便广播器无锁判断用户名是否存在
	checkUserCanInChannel chan bool   //用来回传用户名是否存在

	//获取用户列表
	requestUserListChannel chan struct{} //负责接受请求信号
	userListChannel        chan []*User
}

// BroadCaster 实例化一个广播器，提供给外界使用
var BroadCaster = broadCaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, global.MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

// Start 启动广播器
func (b *broadCaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			//新用户进入
			b.users[user.NickName] = user
			//OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			//用户离开
			delete(b.users, user.NickName)
			//避免goroutine泄露
			user.CloseMessageChanel()

		case msg := <-b.messageChannel:
			//给所有用户发送消息
			for _, user := range b.users {
				//排除自己
				if msg.User.UID == user.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			//OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			_, ok := b.users[nickname]
			if ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUserListChannel:
			users := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				users = append(users, user)
			}
			b.userListChannel <- users
		}
	}

}
func (b *broadCaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

// BroadCast 用于广播信息
func (b *broadCaster) BroadCast(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("message queue is full")
	}
	b.messageChannel <- msg
}
func (b *broadCaster) UserEntering(user *User) {
	b.enteringChannel <- user
}
func (b *broadCaster) UserLeaving(user *User) {
	b.leavingChannel <- user
}
func (b *broadCaster) UserList() []*User {
	b.requestUserListChannel <- struct{}{} //发送信号
	return <-b.userListChannel
}
