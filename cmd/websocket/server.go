package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Hello HTTP")
	})
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := websocket.Accept(writer, request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		//defer 关闭，防止出现错误无法关闭连接
		defer conn.Close(websocket.StatusInternalError, "内部出错了。。。")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("客户端发送信息：%v", v)

		err = wsjson.Write(ctx, conn, "Hello websocket client")
		if err != nil {
			log.Println(err)
			return
		}
		conn.Close(websocket.StatusNormalClosure, "")

	})
	log.Fatal(http.ListenAndServe("2001", nil))

}
