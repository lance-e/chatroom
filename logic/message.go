package logic

import (
	"github.com/spf13/cast"
	"time"
)

const (
	MsgTypeNormal = iota
	MsgTypeWelcome
	MsgTypeEnter
	MsgTypeLeave
	MsgTypeError
)

// Message 消息属性
type Message struct {
	//哪个用户发的
	User           *User     `json:"user"`
	Type           int       `json:"type"`
	Content        string    `json:"content"`
	MsgTime        time.Time `json:"msg_time"`
	ClientSendTime time.Time `json:"client_send_time"`
	//消息@了谁
	Ats []string `json:"ats"`
}

func NewMessage(user *User, content string, clientTime string) *Message {
	msg := &Message{
		User:    user,
		Content: content,
		Type:    MsgTypeNormal,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		msg.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return msg
}

// NewWelcomeMessage 用来给新用户发送欢迎信息
func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: "欢迎来到聊天室:" + user.NickName,
		MsgTime: time.Now(),
	}
}

// NewUserEnterMessage 发送用户进入聊天室的信息
func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeEnter,
		Content: user.NickName + "进入了聊天室",
		MsgTime: time.Now(),
	}
}

// NewUserLeaveMessage 发送用户离开聊天室的信息
func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeLeave,
		Content: user.NickName + "离开了聊天室",
		MsgTime: time.Now(),
	}
}

// NewErrorMessage 返回错误信息
func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
