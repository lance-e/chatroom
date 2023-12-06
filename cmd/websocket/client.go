package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws://127.0.0.1:2001/ws", nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "内部错误")

	err = wsjson.Write(ctx, conn, "hello server!you are good")
	if err != nil {
		log.Println(err)
		return
	}

	var v interface{}
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("接收到服务器响应：%v", v)

	conn.Close(websocket.StatusNormalClosure, "")

}
