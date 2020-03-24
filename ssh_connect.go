package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type SSHConnect struct {
	session    *ssh.Session
	stdinPipe  io.WriteCloser
	stdoutPipe io.Reader
	//stdout     *bytes.Buffer
	//stderr     *bytes.Buffer
}

//
func NewSSHConnect(client *ssh.Client) (sshConn *SSHConnect, err error) {
	var (
		session *ssh.Session
		//stdout, stderr *bytes.Buffer
	)
	if session, err = client.NewSession(); err != nil {
		return
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("linux", 80, 40, modes); err != nil {
		return
	}

	/*stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)*/

	pipe, _ := session.StdinPipe()
	stdoutPipe, _ := session.StdoutPipe()

	/*session.Stdout = stdout
	session.Stderr = stderr*/

	if err = session.Shell(); err != nil {
		return
	}

	return &SSHConnect{
		session:    session,
		stdinPipe:  pipe,
		stdoutPipe: stdoutPipe,
		/*stdout:    stdout,
		stderr:    stderr,*/
	}, nil
}

//Receive messages from websocket
func (s *SSHConnect) Recv(conn *websocket.Conn, quit chan int) {
	defer Quit(quit)
	var (
		bytes []byte
		err   error
	)
	for {
		if bytes, err = WsRecv(conn); err != nil {
			return
		}
		if len(bytes) > 0 {
			if _, e := s.stdinPipe.Write(bytes); e != nil {
				return
			}
		}
	}
}

func (s *SSHConnect) Output(conn *websocket.Conn, quit chan int) {
	defer Quit(quit)
	var (
		read int
		err  error
	)
	tick := time.NewTicker(60 * time.Millisecond)
	defer tick.Stop()
Loop:
	for {
		select {
		case <-tick.C:
			i := make([]byte, 1024)
			if read, err = s.stdoutPipe.Read(i); err != nil {
				fmt.Println(err)
				break Loop
			}
			if err = WsSendText(conn, i[:read]); err != nil {
				fmt.Println(err)
				break Loop
			}
		}
	}
}

//test
func (s *SSHConnect) recvv(command string) {
	if _, err := s.stdinPipe.Write([]byte(command)); err != nil {
		fmt.Println(err)
	}
}

//test
func (s *SSHConnect) output() {
	tick := time.NewTicker(120 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			i := make([]byte, 1024)
			if read, err := s.stdoutPipe.Read(i); err == nil {
				i2 := string(i[:read])
				//Get head
				//split := strings.Split( i2,"\n")
				//fmt.Println(split[len(split)-1])
				fmt.Println(i2)
			}
		}
	}
}

func Quit(quit chan int) {
	quit <- 1
}
