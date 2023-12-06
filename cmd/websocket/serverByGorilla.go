package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(writer, request, nil)
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		var v interface{}
		err = conn.ReadJSON(&v)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("客户端发送的消息：%v", v)

		err = conn.WriteJSON("hello websocket client")
		if err != nil {
			log.Println(err)
			return
		}

	})
	log.Fatal(http.ListenAndServe(":2001", nil))

}
