package attack

import (
	"fmt"
	"github.com/cbot/targets/source"
	"io"
	"os"
	"strings"
)

type AttackDump struct {
	attackTasks *AttackTasks
	out io.Writer
}

func NewAttackDump(attackTasks *AttackTasks,fpath string) *AttackDump {

	out := os.Stdout

	if fpath != "" {

		file, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)

		if err == nil {

			out = file
		}
	}

	return &AttackDump{
		attackTasks: attackTasks,
		out:out,
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

	defer ad.attackTasks.PubUnSyn()

	fmt.Fprintf(ad.out,"%s:%d\n", target.IP(),target.Port())

	return nil
}

func (ad *AttackDump) PubProcess(process *AttackProcess) {

}