package node

import (
	"github.com/cbot/attack"
	"github.com/cbot/attack/ascript"
	"github.com/cbot/targets/source"
)

type Node struct {

	cfg *Config

	spool *source.SourcePool

	attackTasks *attack.AttackTasks


}


func NewNode(cfg *Config) *Node {

	return &Node{
		cfg:         cfg,
		spool:       nil,
		attackTasks: nil,
	}

}

func (n *Node) Start() error {


	return nil
}


func (n *Node) NodeId() string {

	return n.cfg.NodeId
}


func (n *Node) AddAttackSource(name string,types []string,content []byte) error {

	s,err:= source.NewScriptSourceFromContent(n.spool,name,types,content)

	if err!=nil {

		return err
	}

	n.spool.StartSource(s)

	return nil
}

func (n *Node) AddAttack(name string,attackType string,defaultProto string,defaultPort int,content []byte) error {


	att,err := ascript.NewAttackScriptFromContent(n.attackTasks,name,attackType,defaultPort,defaultProto,content)

	if err!=nil {

		return err
	}

	n.attackTasks.AddAttack(att)

	return nil
}