package logic

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"regexp"
	"time"
)

type User struct {
	UID           int           `json:"uid,omitempty"`
	NickName      string        `json:"nickname,omitempty"`
	EnterAt       time.Time     `json:"enter_at"`
	Addr          string        `json:"addr,omitempty"`
	MessageChanel chan *Message `json:"_"`

	//token 进行用户校验
	Token string `json:"token"`

	//isBool 判断用户是否是第一次进入聊天室
	IsBool bool `json:"is_bool"`

	Conn *websocket.Conn
}

// 系统用户 ， 代表系统发送的信息
var System = &User{}

// NewUser 新建一个user实例
func NewUser(conn *websocket.Conn, nickname string, token string, addr string) *User {
	newUser := &User{
		NickName:      nickname,
		EnterAt:       time.Now(),
		Token:         token,
		Addr:          addr,
		MessageChanel: make(chan *Message, 8),
		Conn:          conn,
	}
	if newUser.Token != "" {

	}
	if newUser.UID == 0 {

	}
}
func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChanel {
		u.Conn.WriteJSON(msg)
	}
}
func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = u.Conn.ReadJSON(&receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}
		//内容发送到聊天室
		msg := NewMessage(u, receiveMsg["content"])

		//解析 content ，看看@了谁
		reg := regexp.MustCompile(`@[^\s@]{4,20}`) //?????????
		msg.Ats = reg.FindAllString(msg.Content, -1)

		BroadCaster.BroadCast(msg)
	}

}
