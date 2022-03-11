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

// RFC 4254 Section 6.5.
type execMsg struct {
	Command string
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


func makeCmd(cmd string) []byte {

	req := execMsg{Command:cmd}

	return ssh.Marshal(&req)
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

func (c *SSHClient) SendCmd(cmd string) error {

	//创建ssh-session
	session, err := c.con.NewSession()
	if err != nil {

		return err
	}

	defer session.Close()

	ok, err := session.SendRequest("exec",false,makeCmd(cmd))

	if err == nil && !ok {
		err = fmt.Errorf("ssh: command %s failed", cmd)
	}

	if err != nil {
		return err
	}

	return nil
}
