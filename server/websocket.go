package server

import (
	"chatroom/logic"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func websocketHandleFunc(writer http.ResponseWriter, request *http.Request) {

	//新用户进入，在协议升级之前，把用户名拿到
	nickname := request.FormValue("nickname")

	//将连接升级为websocket连接，
	ws := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return false
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := ws.Upgrade(writer, request, nil)
	if err != nil {
		log.Printf("websocket upgrade error:%v", err)
		return
	}
	//防止中途出现问题，导致无法正常关闭
	defer conn.Close()

	//对昵称长度进行判断
	if len(nickname) < 4 || len(nickname) > 20 {
		log.Println("nickname illegal :", nickname)
		conn.WriteJSON("非法昵称，昵称长度应为4-20")
		conn.Close()
		return
	}

	if !logic.BroadCaster.CanEnterRoom(nickname) {
		log.Printf("昵称 : %s 已存在", nickname)
		conn.WriteJSON("该昵称已存在，请更换昵称")
		conn.Close()
		return
	}
	//创建用户实例
	user := logic.NewUser(conn, nickname, request.RemoteAddr)

	//开启给用户发送消息的协程
	go user.SendMessage(request.Context())

	//给新用户发送欢迎信息
	user.MessageChanel <- logic.NewWelcomeMessage(nickname)
	//再给所有用户发送新用户进入聊天室的信息
	msg := logic.NewNoticeMessage(nickname + "进入聊天室")
	logic.BroadCaster.BroadCast(msg)

	//将用户加入广播器的用户列表中
	logic.BroadCaster.UserEntering(user)
	log.Println("user: `" + nickname + "` joins chat")

	//接收用户消息
	err = user.ReceiveMessage(request.Context())
	// 用户离开,给所有用户发送用户离开聊天室的信息
	logic.BroadCaster.UserLeaving(user)
	msg = logic.NewNoticeMessage(nickname + "离开聊天室")
	logic.BroadCaster.BroadCast(msg)
	log.Println("user: `" + nickname + "` leaves chat")
	//关闭
	if err == nil {
		log.Println("connection close...")
		conn.Close()
	} else {
		log.Println("read from client error : ", err)
		conn.Close()
	}
}
