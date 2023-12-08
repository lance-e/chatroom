package logic

// broadCast 广播器
type broadCaster struct {
	//所有聊天室用户
	users map[string]*User

	//所有goroutine统一管理，避免外部乱用

	enteringChanel chan *User
	leavingChanel  chan *User
	messageChanel  chan *Message

	checkUserChanel      chan string //接收昵称，方便广播器无锁判断用户名是否存在
	checkUserCanInChanel chan bool   //用来回传用户名是否存在
}

var MessageQueueLen = 10

// BroadCaster 实例化一个广播器，提供给外界使用
var BroadCaster = broadCaster{
	users: make(map[string]*User),

	enteringChanel: make(chan *User),
	leavingChanel:  make(chan *User),
	messageChanel:  make(chan *Message, MessageQueueLen),

	checkUserChanel:      make(chan string),
	checkUserCanInChanel: make(chan bool),
}

// Start 启动广播器
func (b *broadCaster) Start() {
	for {
		select {
		case user := <-b.enteringChanel:
			//新用户进入
			b.users[user.NickName] = user
			b.sendUserList()
			OfflineProcessor.Send(user)
		case user := <-b.leavingChanel:
			//用户离开
			delete(b.users, user.NickName)
			//避免goroutine泄露
			b.CloseMessageChanel()
			b.sendUserList()
		case msg := <-b.messageChanel:
			//给所有用户发送消息
			for _, user := range b.users {
				//排除自己
				if msg.User.UID == user.UID {
					continue
				}
				user.MessageChanel <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChanel:
			_, ok := b.users[nickname]
			if !ok {
				b.checkUserCanInChanel <- false
			} else {
				b.checkUserCanInChanel <- true
			}
		}
	}

}
func (b *broadCaster) CanEnterRoom(name string) bool {
	b.checkUserChanel <- name
	return <-b.checkUserCanInChanel
}

// BroadCast 用于广播信息
func (b *broadCaster) BroadCast(msg *Message) {
	b.messageChanel <- msg
}
func (b *broadCaster) UserEntering(user *User) {
	b.enteringChanel <- user
}
func (b *broadCaster) UserLeaving(user *User) {
	b.leavingChanel <- user
}
