package logic

// broadCast 广播器
type BroadCast struct {
	//所有聊天室用户
	users map[string]*User

	//所有goroutine统一管理，避免外部乱用

	enteringChanel chan *User
	leavingChanel  chan *User
	messageChanel  chan *User

	checkUserChanel      chan string
	checkUserCanInChanel chan bool
}

func (b *BroadCast) Start() {

}
func (b *BroadCast) CanEnterRoom(name string) bool {
	b.checkUserChanel <- name
	return <-b.checkUserCanInChanel
}
