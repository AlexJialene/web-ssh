package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

//create ssh client
func createSSHClient(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		//session      *ssh.Session
		err error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//Handling the host key
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return client, nil
	/*if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil*/
}

func runSSH(client *ssh.Client, command string) (string, error) {
	var err error
	var session *ssh.Session
	if session, err = client.NewSession(); err == nil {
		session.StdinPipe()
		defer session.Close()
		var stdOut bytes.Buffer

		session.Stdout = &stdOut
		err = session.Run(command)
		if err != nil {
			return "", err
		}

		return string(stdOut.Bytes()), nil
	}
	return "", err

}
