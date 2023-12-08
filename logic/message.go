package logic

import (
	"time"
)

const (
	MsgTypeNormal = iota
	MsgTypeSystem
	MsgTypeError
	MsgTypeUserList
)

// Message 消息属性
type Message struct {
	//哪个用户发的
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	//消息@了谁
	Ats []string `json:"ats"`

	Users map[string]*User
}

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Content: content,
	}
}

// NewWelcomeMessage 用来给新用户发送欢迎信息
func NewWelcomeMessage(nickname string) *Message {
	return &Message{
		//User: &User{
		//	NickName: nickname,
		//},
		Type:    MsgTypeSystem,
		Content: "欢迎来到聊天室:" + nickname,
		MsgTime: time.Now(),
	}
}
func NewNoticeMessage(msg string) {
	//return &Message{
	//	Type:    MsgTypeSystem,
	//	Content: msg,
	//	MsgTime: time.Now(),
	//}
}
