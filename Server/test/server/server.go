package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 将HTTP连接升级成WebSocket连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有跨域连接
		return true
	},
}

// 处理WebSocket请求
func handler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		conn *websocket.Conn
	)

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// 接收客户端消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				return
			}
		}
		log.Printf("recv: %s", message)

		// 发送客户端消息
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println(err)
		}
		log.Printf("sent: %s", message)
	}
}

func main() {
	// 设置WebSocket路由
	http.HandleFunc("/ws", handler)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
