package server

import (
	"chatroom/logic"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func websocketHandleFunc(writer http.ResponseWriter, request *http.Request) {

	//新用户进入，在协议升级之前，把用户信息拿到
	nickname := request.FormValue("nickname")
	token := request.FormValue("token")

	//将连接升级为websocket连接，
	ws := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
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
		conn.WriteJSON(logic.NewErrorMessage("非法昵称，昵称长度应为4-20"))
		conn.Close()
		return
	}

	if !logic.BroadCaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在", nickname)
		conn.WriteJSON(logic.NewErrorMessage("用户名已存在，换一个名字吧"))
		conn.Close()
		return
	}
	//创建用户实例
	userHasToken := logic.NewUser(conn, nickname, token, request.RemoteAddr)

	//开启给用户发送消息的协程
	go userHasToken.SendMessage()

	//给新用户发送欢迎信息
	userHasToken.MessageChannel <- logic.NewWelcomeMessage(userHasToken)

	//避免token泄露，在这里进行token处理
	temUser := *userHasToken
	user := &temUser
	user.Token = ""

	//再给所有用户发送新用户进入聊天室的信息
	msg := logic.NewUserEnterMessage(user)
	logic.BroadCaster.BroadCast(msg)

	//将用户加入广播器的用户列表中
	logic.BroadCaster.UserEntering(user)
	log.Println("user: `" + nickname + "` joins chat")

	//接收用户消息
	err = user.ReceiveMessage()

	// 用户离开,给所有用户发送用户离开聊天室的信息
	logic.BroadCaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.BroadCaster.BroadCast(msg)
	log.Println("user: `" + nickname + "` leaves chat")

	//关闭连接
	if err == nil {
		log.Println("connection close...")
		conn.Close()
	} else {
		log.Println("read from client error : ", err)
		conn.Close()
	}
}
