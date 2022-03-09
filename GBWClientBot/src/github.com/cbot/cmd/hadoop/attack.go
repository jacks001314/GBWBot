package main

import (
	"flag"
	"github.com/cbot/proto/hadoop/hadoop_yarn"
	yarn_conf "github.com/cbot/proto/hadoop/hadoop_yarn/conf"
	"github.com/cbot/proto/hadoop/hadoop_yarn/yarn_client"
	"log"
	"time"
)

func main(){

	addr := flag.String("addr","127.0.0.1:8032","set the hadoop ipc address")
	appName := flag.String("name","sbot_hadoop","set the yarn application name")
	appType := flag.String("type","sbot_hadoop_type","set the yarn application type")
	queue   := flag.String("queue","default","set the yarn application queue name")

	cmd := flag.String("cmd","touch /var/tmp/fuck.data","set the command that to been run")

	flag.Parse()

	// Create YarnConfiguration
	conf, _ := yarn_conf.NewYarnConfiguration()

	conf.SetRMAddress(*addr)

	// Create YarnClient
	yarnClient, err := yarn_client.CreateYarnClient(conf)
	if err!=nil {
		log.Fatalf("Create hadoop yarn client failed:%v",err)
	}

	// Create new application to get ApplicationSubmissionContext
	_, asc, err := yarnClient.CreateNewApplication()

	if err!=nil {
		log.Fatalf("create a new hadoop yarn application failed:%v",err)

	}

	log.Printf("Create a new hadoop yarn application ok,appID:%v\n",asc.ApplicationId)

	// Setup ContainerLaunchContext for the application
	clc := hadoop_yarn.ContainerLaunchContextProto{}
	clc.Command = []string{*cmd}

	// Resource for ApplicationMaster
	var memory int32 = 1024
	resource := hadoop_yarn.ResourceProto{Memory: &memory}

	// Setup ApplicationSubmissionContext for the application
	asc.AmContainerSpec = &clc
	asc.Resource = &resource
	asc.ApplicationName = appName
	asc.Queue = queue
	asc.ApplicationType = appType

	// Submit!
	err = yarnClient.SubmitApplication(asc)
	if err != nil {
		log.Fatal("yarnClient.SubmitApplication ", err)
	}

	log.Println("Successfully submitted application: ", asc.ApplicationId)

	appReport, err := yarnClient.GetApplicationReport(asc.ApplicationId)
	if err != nil {
		log.Fatal("yarnClient.GetApplicationReport ", err)
	}
	appState := appReport.GetYarnApplicationState()
	for appState != hadoop_yarn.YarnApplicationStateProto_FINISHED && appState != hadoop_yarn.YarnApplicationStateProto_KILLED && appState != hadoop_yarn.YarnApplicationStateProto_FAILED {
		log.Println("Application in state ", appState)
		time.Sleep(1 * time.Second)
		appReport, err = yarnClient.GetApplicationReport(asc.ApplicationId)
		appState = appReport.GetYarnApplicationState()
	}

	log.Println("Application finished in state: ", appState)
}



