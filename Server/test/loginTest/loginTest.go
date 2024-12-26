package main

import (
	"Common/Framework/codec"
	"Common/message"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"reflect"
)

func main() {
	var (
		err error
	)

	err = codec.Init()
	if err != nil {
		fmt.Println(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9000/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer conn.Close()

	sendPrelogin(conn)
	onRead(conn)
}

func sendPrelogin(conn *websocket.Conn) {
	preloginReq := message.PreLoginReq{CheckCode: 0x32388545}
	buf, err := codec.Encode(preloginReq)
	if err != nil {
		log.Fatal("encode:", err)
	}

	err = conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		log.Fatal("write:", err)
	}
}

func onRead(conn *websocket.Conn) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}

		if buf != nil {
			message, err := codec.Decode(buf)
			if err != nil {
				log.Fatal("decode:", err)
				return
			}
			log.Println("receive type:", reflect.TypeOf(message), message)
			onMessage(conn, message)
		}
	}
}

func onMessage(conn *websocket.Conn, info interface{}) {
	var (
		err error
		buf []byte
	)
	switch info.(type) {
	case message.PreLoginRes:
		// todo 填写内容
		msg := message.LoginReq{
			Account: "guanhui",
		}

		buf, err = codec.Encode(msg)
		if err != nil {
			log.Println("encode:", err)
		}
	case message.LoginRes:
		msg := message.EnterReq{
			Zone: 1,
		}

		buf, err = codec.Encode(msg)
		if err != nil {
			log.Println("encode:", err)
		}
	case message.EnterRes:
		fmt.Println("enter success!")
	}

	if buf != nil {
		err = conn.WriteMessage(websocket.BinaryMessage, buf)
		if err != nil {
			log.Println("write:", err)
		}
	}
}
