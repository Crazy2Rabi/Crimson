package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	var err error

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9000/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer conn.Close()

	// 发送消息
	err = conn.WriteMessage(websocket.BinaryMessage, []byte("hello"))
	if err != nil {
		log.Fatal("write:", err)
		return
	}

	// 接收消息
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
		return
	}
	fmt.Println(string(message))

	// 关闭
	err = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatal("write close:", err)
		return
	}
}
