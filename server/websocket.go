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
	defer conn.Close()

	//对昵称长度进行判断
	if len(nickname) < 4 || len(nickname) > 20 {
		log.Println("nickname illegal :", nickname)
		conn.WriteJSON("非法昵称，昵称长度应为4-20")
		conn.Close()
		return
	}
	var broadcaster = logic.BroadCast{}
	if !broadcaster.CanEnterRoom(nickname) {
		log.Printf("昵称 : %s 已存在", nickname)
		conn.WriteJSON("该昵称已存在，请更换昵称")
		conn.Close()
		return
	}
	//创建用户实例
	user := logic.NewUser(conn, nickname, request.RemoteAddr)

	//开启给用户发送消息的协程
	go user.SendMessage(request.Context())
}
