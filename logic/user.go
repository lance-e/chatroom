package logic

import (
	"context"
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

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	return &User{
		Conn:     conn,
		NickName: nickname,
		Addr:     addr,
	}
}
func (u *User) SendMessage(ctx context.Context) {

}
