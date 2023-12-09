package server

import (
	"chatroom/logic"
	"net/http"
)

func RegisterHandle() {

	//接收广播消息
	go logic.BroadCaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/user_list", UserListHandleFunc)
	http.HandleFunc("/ws", websocketHandleFunc)
}
