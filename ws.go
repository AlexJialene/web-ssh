package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Command string
}

func WsRecv(conn *websocket.Conn) ([]byte, error) {
	var (
		err  error
		data []byte
	)
	_, data, err = conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	return data, nil

}

func WsSendText(conn *websocket.Conn, b []byte) error {
	if err := conn.WriteMessage(1, b); err != nil {
		return err
	}
	return nil

}
