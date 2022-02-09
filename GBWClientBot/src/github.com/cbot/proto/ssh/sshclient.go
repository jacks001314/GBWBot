package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"time"
)

type SSHClient struct {

	con *ssh.Client
}

func LoginWithPasswd(host string,port int,user string,passwd string,timeout int64) (*SSHClient,error) {

	config := &ssh.ClientConfig{
		User:              user,
		Auth:              []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		Timeout:           time.Duration(timeout)*time.Millisecond,
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host, port)

	con, err := ssh.Dial("tcp", addr, config)

	if err != nil {

		return nil,err
	}

	return &SSHClient{con:con},nil
}

func LoginWithPrivKey(host string,port int,user string,privkey string,timeout int64) (*SSHClient,error) {

	signer, err := ssh.ParsePrivateKey([]byte(privkey))
	if err != nil {

		return nil,err
	}

	clientkey := ssh.PublicKeys(signer)

	config := &ssh.ClientConfig{
		User:              user,
		Auth:              []ssh.AuthMethod{clientkey},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		Timeout:           time.Duration(timeout)*time.Millisecond,
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host, port)

	con, err := ssh.Dial("tcp", addr, config)

	if err != nil {

		return nil,err
	}

	return &SSHClient{con:con},nil
}

func LoginNoPassword(host string,port int,user string,timeout int64) (*SSHClient,error) {

	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil,err
	}

	agent := agent.NewClient(sock)

	signers, err := agent.Signers()
	if err != nil {
		return nil,err
	}

	auths := []ssh.AuthMethod{ssh.PublicKeys(signers...)}

	config := &ssh.ClientConfig{
		User:              user,
		Auth:              auths,
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		Timeout:           time.Duration(timeout)*time.Millisecond,
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host, port)

	con, err := ssh.Dial("tcp", addr, config)

	if err != nil {

		return nil,err
	}

	return &SSHClient{con:con},nil
}


func (c *SSHClient) Close(){

	c.con.Close()
}

func (c *SSHClient) RunCmd(cmd string) ([]byte,error){

	//创建ssh-session
	session, err := c.con.NewSession()
	if err != nil {

		return nil,err
	}

	defer session.Close()

	return session.CombinedOutput(cmd)
}

