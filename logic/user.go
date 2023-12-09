package logic

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io"
	"regexp"
	"strings"
	"time"
)

var globalUID uint32 = 0

type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	//token 进行用户校验
	Token string `json:"token"`

	//isBool 判断用户是否是第一次进入聊天室
	isNew bool

	conn *websocket.Conn
}

// 系统用户 ， 代表系统发送的信息
var System = &User{}

// NewUser 新建一个user实例
func NewUser(conn *websocket.Conn, nickname string, token string, addr string) *User {
	user := &User{
		NickName:       nickname,
		EnterAt:        time.Now(),
		Token:          token,
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		conn:           conn,
	}
	if user.Token != "" {
		uid, err := parseTokenAndValid(user.Token, user.NickName)
		if err == nil {
			user.UID = uid
		}
	}
	if user.UID == 0 {
		//使用uuid库来生成uid吧
		UUID, _ := uuid.NewUUID()
		user.UID = int(UUID.ID())

		//user.UID = int(atomic.AddUint32(&globalUID, 1))
		user.Token = genToken(user.UID, user.NickName)
		user.isNew = true

	}
	return user
}
func (u *User) SendMessage() {
	for msg := range u.MessageChannel {
		_ = u.conn.WriteJSON(msg)
	}
}
func (u *User) ReceiveMessage() error {
	var receiveMsg map[string]string
	var err error
	for {
		err = u.conn.ReadJSON(&receiveMsg)

		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		//内容发送到聊天室
		msg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		msg.Content = FilterSensitive(msg.Content) //敏感词汇过滤
		//解析 content ，看看@了谁
		reg := regexp.MustCompile(`@[^\s@]{4,20}`) //?????????
		msg.Ats = reg.FindAllString(msg.Content, -1)

		BroadCaster.BroadCast(msg)
	}

}

// CloseMessageChanel 用于关闭消息通道,避免goroutine泄露
func (u *User) CloseMessageChanel() {
	close(u.MessageChannel)
}

// genToken 生成token
func genToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")               //获取密钥
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid) //先将nickname，secret，uid拼接

	messageMac := macSHA256([]byte(message), []byte(secret)) //使用hmac-SHA256加密

	token := fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMac), uid) //再将加密后的token再拼接uid
	return token

}

// macSHA256 hmac-SHA256加密
func macSHA256(msg, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(msg)
	return mac.Sum(nil)
}

// parseTokenAndValid 解析token内容，分析是否有效,返回UID和error
func parseTokenAndValid(token string, nickname string) (int, error) {
	//获取token中的信息

	index := strings.LastIndex(token, "uid")
	messageMac, err := base64.StdEncoding.DecodeString(token[:index]) //对genToken中的messageMac进行解码
	if err != nil {
		return 0, err
	}

	uid := cast.ToInt(token[index+3:]) // 获取token中的uid信息

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid) //先将nickname，secret，uid拼接

	//进行token校验
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	hash := mac.Sum(nil)
	if ok := hmac.Equal(messageMac, hash); ok {
		return uid, nil
	}
	return 0, errors.New("this token is illegal")
}
