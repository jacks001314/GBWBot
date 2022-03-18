package handler

import (
	"container/list"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/utils/fileutils"
	"github.com/sbot/utils/jsonutils"
	"io/ioutil"
	"sync"
	"time"
)

type AttackScriptsHandler struct {

	lock sync.Mutex
	attackScriptsConfig *AttackScriptsConfig

	queueCapacity uint32

	requests *list.List

	waitingRequests uint32

	waitFetchTimeout uint64
}


type AttackScriptsConfig struct {

	StorePath string `json:"storePath"`
	StoreRPath string `json:"storeRPath"`
	AttackScripts []*AttackScriptEntry `json:"attackScripts"`
}


type AttackScriptEntry struct {

	Id string `json:"id"`
	Name string `json:"name"`
	DefaultProto string `json:"defaultProto"`
	DefaultPort int `json:"defaultPort"`
	AttackType string `json:"attackType"`
	Fname string `json:"fname"`
	Inline bool `json:"inline"`
	Language string `json:"language"`
	Desc string `json:"desc"`

}

type attackScriptsContext struct {

	scripts []*model.AttackScripts

	listNode *list.Element

	nodesMap map[string]bool

	waitedFetchs int

	waitChan chan bool
}

func NewAttackScriptsHandler(cfgPath string,queueCapacity uint32,waitFetchTimeout uint64) (*AttackScriptsHandler,error) {

	var cfg AttackScriptsConfig

	if err := jsonutils.UNMarshalFromFile(&cfg,cfgPath);err!=nil {

		errS := fmt.Sprintf("Cannot load attack Scripts config from file:%s,failed:%v",cfgPath,err)
		log.Errorf(errS)

		return nil,fmt.Errorf(errS)
	}

	return &AttackScriptsHandler {
		lock:                sync.Mutex{},
		attackScriptsConfig: &cfg,
		queueCapacity:       queueCapacity,
		requests:            list.New(),
		waitingRequests:     0,
		waitFetchTimeout:    waitFetchTimeout,
	},nil
}


func (ash* AttackScriptsHandler) findAttackScriptsEntry(cfgId string) *AttackScriptEntry {

	for _,entry:= range ash.attackScriptsConfig.AttackScripts {

		if entry.Id == cfgId {
			return entry
		}
	}

	return nil
}

func (ash* AttackScriptsHandler) makeAttackScriptsContent(entry *AttackScriptEntry) ([]byte,error) {

	fpath := fmt.Sprintf("%s%s%s",ash.attackScriptsConfig.StorePath,ash.attackScriptsConfig.StoreRPath,entry.Fname)

	if !fileutils.FileIsExisted(fpath){

		errS := fmt.Sprintf("The Attack Scripts file:%s not existed!",fpath)
		log.Errorf(errS)

		return nil,fmt.Errorf(errS)
	}

	return ioutil.ReadFile(fpath)
}

func (ash* AttackScriptsHandler) addAttackScriptsContext(ctx *attackScriptsContext) {

	ash.lock.Lock()
	defer ash.lock.Unlock()
	ctx.listNode = ash.requests.PushBack(ctx)

	ash.waitingRequests = ash.waitingRequests+1

}

func (ash* AttackScriptsHandler) removeAttackScriptsContext(ctx *attackScriptsContext) {

	ash.lock.Lock()
	defer ash.lock.Unlock()

	ash.requests.Remove(ctx.listNode)
	ash.waitingRequests = ash.waitingRequests-1

}

func (ash* AttackScriptsHandler) findAttackScriptsContext(nodeId string) *attackScriptsContext {

	ash.lock.Lock()
	defer ash.lock.Unlock()

	for listNode := ash.requests.Front();listNode!=nil;listNode = listNode.Next() {

		ctx := listNode.Value.(*attackScriptsContext)

		if fetched,ok := ctx.nodesMap[nodeId];ok&&!fetched {

			return ctx
		}
	}

	//not existed
	return nil
}

func (ash* AttackScriptsHandler) fetchedAttackScriptsContext(nodeId string,ctx *attackScriptsContext){

	ash.lock.Lock()
	defer ash.lock.Unlock()

	ctx.nodesMap[nodeId] = true

	ctx.waitedFetchs = ctx.waitedFetchs-1

	if ctx.waitedFetchs <= 0 {

		ctx.waitChan<-true

	}
}

func (ash* AttackScriptsHandler) makeAttackScriptsContext(request *model.AddAttackScriptsRequest) *attackScriptsContext {

	nodesMap := make(map[string]bool)

	for _,nodeId := range request.NodeIds {
		nodesMap[nodeId] = false
	}

	ctx := &attackScriptsContext{
		scripts:      make([]*model.AttackScripts,0),
		listNode:     nil ,
		nodesMap:     nodesMap,
		waitedFetchs: len(request.NodeIds),
		waitChan:     make(chan bool),
	}

	for _,cfgId := range request.CfgIds {

		entry := ash.findAttackScriptsEntry(cfgId)

		if entry !=nil &&!entry.Inline {

			content,err := ash.makeAttackScriptsContent(entry)

			if err!=nil {
				log.Errorf(err.Error())
				continue
			}

			ctx.scripts= append(ctx.scripts, &model.AttackScripts{
				Name:         entry.Name,
				AttackType:   entry.AttackType,
				DefaultPort:  int32(entry.DefaultPort),
				DefaultProto: entry.DefaultProto,
				Size:        uint64(len(content)),
				Content:      content,
				HasNext:      true,
			})
		}
	}

	if len(ctx.scripts)>0 {

		ctx.scripts[len(ctx.scripts)-1].HasNext = false
		ash.addAttackScriptsContext(ctx)
		return ctx
	}

	return nil
}

func (ash* AttackScriptsHandler) waitAttackScriptsFetched(requestCtx *attackScriptsContext) *model.AddAttackScriptsReply{

	var msg string
	tm := time.NewTicker(time.Duration(ash.waitFetchTimeout)*time.Millisecond)

	defer tm.Stop()

	for {
		select {
		case <-tm.C:
			//timeout
			msg = "timeout"
			goto out

		case <-requestCtx.waitChan:
			msg = "ok"
			goto out
		}
	}

out:
	fetchedNodes := make([]string,0)
	for id,fetched := range requestCtx.nodesMap {

		if fetched {
			fetchedNodes = append(fetchedNodes,id)
		}
	}

	ash.removeAttackScriptsContext(requestCtx)

	return &model.AddAttackScriptsReply{
		Status:       0,
		Message:      msg,
		FetchedNodes: fetchedNodes,
	}
}

func (ash* AttackScriptsHandler) AddAttackScriptsHandle(request *model.AddAttackScriptsRequest) *model.AddAttackScriptsReply {

	log.Infof("Accept an add attack scripts request:%v",request)

	if request.NodeIds == nil || len(request.NodeIds) == 0 {

		errs := fmt.Sprintf("Please Specify the nodes that accept this attack scripts:%d",ash.waitingRequests)
		log.Errorf(errs)

		return &model.AddAttackScriptsReply{
			Status:       -1,
			Message:      errs,
			FetchedNodes: []string{},
		}
	}

	if ash.waitingRequests+1>ash.queueCapacity {

		errs := fmt.Sprintf("Too many add attack scripts request to waiting,current requests:%d",ash.waitingRequests)
		log.Errorf(errs)

		return &model.AddAttackScriptsReply{
			Status:       -1,
			Message:      errs,
			FetchedNodes: []string{},
		}
	}


	ctx := ash.makeAttackScriptsContext(request)

	if ctx == nil  {


		errs := fmt.Sprintf("no attack scripts cann been created,please specify valid config ids:%v",request.CfgIds)
		log.Errorf(errs)

		return &model.AddAttackScriptsReply{
			Status:       -1,
			Message:      errs,
			FetchedNodes: []string{},
		}
	}

	//wait request to be accept by all nodes

	return ash.waitAttackScriptsFetched(ctx)

}

func (ash *AttackScriptsHandler) FetchAttackScripts(request *model.FetchAttackScriptsRequest) []*model.AttackScripts {


	ctx := ash.findAttackScriptsContext(request.NodeId)

	if ctx == nil {

		//log.Infof("no find some attack scripts need to be added for node:%s",request.NodeId)

		return []*model.AttackScripts{}
	}

	log.Infof("find some attack scripts need to be added for node:%s,scripts:%d",request.NodeId,len(ctx.scripts))

	ash.fetchedAttackScriptsContext(request.NodeId,ctx)

	return ctx.scripts
}
