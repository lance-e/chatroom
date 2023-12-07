package logic

// broadCast 广播器
type broadCaster struct {
	//所有聊天室用户
	users map[string]*User

	//所有goroutine统一管理，避免外部乱用

	enteringChanel chan *User
	leavingChanel  chan *User
	messageChanel  chan *User

	checkUserChanel      chan string //接收昵称，方便广播器无锁判断用户名是否存在
	checkUserCanInChanel chan bool   //用来回传用户名是否存在
}

// BroadCaster 实例化一个广播器，提供给外界使用
var BroadCaster = broadCaster{
	users: make(map[string]*User),

	enteringChanel: make(chan *User),
	leavingChanel:  make(chan *User),
	messageChanel:  make(chan *User),

	checkUserChanel:      make(chan string),
	checkUserCanInChanel: make(chan bool),
}

func (b *broadCaster) Start() {

}
func (b *broadCaster) CanEnterRoom(name string) bool {
	b.checkUserChanel <- name
	return <-b.checkUserCanInChanel
}

// BroadCast 用于广播信息
func (b *broadCaster) BroadCast(msg string) {

}
func (b *broadCaster) UserEntering(user *User) {
	b.users[user.NickName] = user
}
func (b *broadCaster) UserLeaving(user *User) {

}
