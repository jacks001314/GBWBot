package attack

import (
	"github.com/cbot/targets/source"
	"log"
	"strings"
)

type AttackDump struct {
	attackTasks *AttackTasks
}

func NewAttackDump(attackTasks *AttackTasks) *AttackDump {

	return &AttackDump{
		attackTasks: attackTasks,
	}
}

func (ad *AttackDump) Name() string {

	return "attackDump"
}

func (ad *AttackDump) DefaultPort() int {

	return 0
}

func (ad *AttackDump) DefaultProto() string {

	return "dump"
}

func (ad *AttackDump) Accept(target source.Target) bool {

	types := target.Source().GetTypes()

	for _, t := range types {

		if strings.EqualFold(t, "dump") {

			return true
		}
	}

	return false
}


func (ad *AttackDump) Run(target source.Target) error {

	log.Printf("Attack Dump tryto attack targets:{ip:%s,host:%s,port:%d,app:%s,proto:%s}",
		target.IP(),target.Host(),target.Port(),target.App(),target.Proto())

	return nil
}

func (ad *AttackDump) PubProcess(process *AttackProcess) {

}