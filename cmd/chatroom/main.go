package main

import (
	"chatroom/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		addr   = ":2022"
		banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |
	
	聊天室
`
	)
	fmt.Printf(banner+"start on %s", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe("localhost"+addr, nil))
}
