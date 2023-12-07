package logic

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

type User struct {
	UID           int           `json:"uid,omitempty"`
	NickName      string        `json:"nickname,omitempty"`
	EnterAt       time.Time     `json:"enter_at"`
	Addr          string        `json:"addr,omitempty"`
	MessageChanel chan *Message `json:"_"`

	Conn *websocket.Conn
}

// 系统用户 ， 代表系统发送的信息
var System = &User{}

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	return &User{
		Conn:     conn,
		NickName: nickname,
		Addr:     addr,
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
		BroadCaster.BroadCast(msg)
	}

}
