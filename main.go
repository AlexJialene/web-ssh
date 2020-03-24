package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"html/template"
	"net/http"
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
	var (
		conn    *websocket.Conn
		client  *ssh.Client
		sshConn *SSHConnect
		err     error
	)
	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	defer conn.Close()

	//Create ssh client
	if client, err = createSSHClient(user, password, host, port); err != nil {
		WsSendText(conn, []byte(err.Error()))
		return
	}
	defer client.Close()

	//connect to ssh
	if sshConn, err = NewSSHConnect(client); err != nil {
		WsSendText(conn, []byte(err.Error()))
		return
	}

	quit := make(chan int)
	go sshConn.Output(conn, quit)
	go sshConn.Recv(conn, quit)
	<-quit
}

func home(w http.ResponseWriter, r *http.Request) {
	temp, e := template.ParseFiles("./template/index.html")
	if e != nil {
		fmt.Println(e)
	}
	temp.Execute(w, nil)
	return
}

func main() {
	http.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css/"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("static/js/"))))

	http.HandleFunc("/index", home)
	http.HandleFunc("/ws/v1", wsHandle)
	http.ListenAndServe(":8080", nil)
}
