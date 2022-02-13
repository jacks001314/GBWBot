package attack

import (
	"github.com/cbot/targets"
	"github.com/cbot/targets/source"
	"sync"
)

type AttackTasks struct {

	lock sync.Mutex

	cfg *Config

	spool *source.SourcePool

	attacks map[string]Attack

	syncChan chan int

	attackProcessChan chan *AttackProcess

}

func NewAttackTasks(cfg *Config,spool *source.SourcePool) *AttackTasks {


	return &AttackTasks{
		lock:    sync.Mutex{},
		cfg:     cfg,
		spool:   spool,
		attacks: make(map[string]Attack),
		syncChan: make(chan int,cfg.MaxThreads),
		attackProcessChan:  make(chan *AttackProcess,cfg.AttackProcessCapacity),
	}

}

func (at *AttackTasks) AddAttack(attack Attack) {

	at.lock.Lock()
	defer at.lock.Unlock()

	if _,ok := at.attacks[attack.Name()];!ok {

		//no existed

		at.attacks[attack.Name()] = attack
	}

}

func (at *AttackTasks) RemoveAttack(name string) {

	at.lock.Lock()
	defer at.lock.Unlock()

	delete(at.attacks,name)

}

func (at *AttackTasks) SubAttackProcess() chan *AttackProcess {

	return at.attackProcessChan
}

func (at *AttackTasks) PubAttackProcess(process *AttackProcess) {

	at.attackProcessChan<- process

}


func (at *AttackTasks) run(target targets.Target) {

	at.lock.Lock()
	defer at.lock.Unlock()

	for _,attack := range at.attacks {

		if attack.Accept(target) {

			go attack.Run(target)
		}
	}
}


func (at *AttackTasks) PubSyn(){

	at.syncChan<-1

}

func (at *AttackTasks) Start() {

	targetChan := at.spool.SubTarget("attack_tasks",at.cfg.SourceCapacity, func(target targets.Target) bool {

		for _,attack := range at.attacks {

			if attack.Accept(target) {

				return true
			}
		}
		return false
	})

	for {

		select {

		case target := <-targetChan:

			//try to run,if too many threads is live than wait some threads exit
			<-at.syncChan
			//ok
			at.run(target)

		}
	}
}

