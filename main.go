package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	//Allow cross-domain
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	user     = "root"
	password = ""
	host     = ""
	port     = 22
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	//Create ssh client
	if client, err := createSSHClient(user, password, host, port); err == nil {
		defer client.Close()

	}

}

func main() {
	//http.HandleFunc("/ws/v1" ,wsHandle)

	client, err := createSSHClient(user, password, host, port)
	if err != nil {
		log.Fatal(err)
	}

	connect, err := NewSSHConnect(client)
	if err != nil {
		log.Fatal(err)
	}

	//connect.recvv("ll")
	go connect.output()
	time.Sleep(5000 * time.Millisecond)
	connect.recvv("ll \n")

	time.Sleep(5000 * time.Millisecond)
	connect.recvv("docker ps \n")

	time.Sleep(5000 * time.Millisecond)
	connect.recvv("cd /opt")

	time.Sleep(5000 * time.Millisecond)
	connect.recvv("ll")

	select {}
}
