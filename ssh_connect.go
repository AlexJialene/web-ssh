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
func (s *SSHConnect) recv(wc *websocket.Conn) {
	//todo
}

func (s *SSHConnect) recvv(command string) {
	s.stdinPipe.Write([]byte(command))
}

func (s *SSHConnect) output() {
	tick := time.NewTicker(120 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			i := make([]byte, 1024)
			if read, err := s.stdoutPipe.Read(i); err == nil {
				fmt.Println(string(i[:read]))
			}
		}
	}
}
