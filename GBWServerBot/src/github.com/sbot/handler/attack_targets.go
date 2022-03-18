package handler

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/utils/fileutils"
	"github.com/sbot/utils/jsonutils"
	"io/ioutil"
	"strings"
	"sync"
	"text/template"
	"time"
)

type AttackTargetsHandler struct {

	lock sync.Mutex
	attackTargetsConfig *AttackTargetsConfig

	queueCapacity uint32

	requests *list.List

	waitingRequests uint32

	waitFetchTimeout uint64
}


type AttackTargetsConfig struct {

	StorePath string `json:"storePath"`
	StoreRPath string `json:"storeRPath"`
	AttackTargets []*AttackTargetEntry `json:"attackTargets"`
}

type AttackTargetEntry struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Fname string `json:"fname"`
	IsTemplate bool `json:"isTemplate"`
	Desc string `json:"desc"`
}

type attackTargetsContext struct {

	targets *model.AttackTargets

	listNode *list.Element

	nodesMap map[string]bool

	waitedFetchs int

	waitChan chan bool
}

func NewAttackTargetsHandler(cfgPath string,queueCapacity uint32,waitFetchTimeout uint64) (*AttackTargetsHandler,error) {

	var cfg AttackTargetsConfig

	if err := jsonutils.UNMarshalFromFile(&cfg,cfgPath);err!=nil {

		errS := fmt.Sprintf("Cannot load attack targets config from file:%s,failed:%v",cfgPath,err)
		log.Errorf(errS)

		return nil,fmt.Errorf(errS)
	}

	return &AttackTargetsHandler{
		lock:                sync.Mutex{},
		attackTargetsConfig: &cfg,
		queueCapacity:       queueCapacity,
		requests:            list.New(),
		waitingRequests:     0,
		waitFetchTimeout:    waitFetchTimeout,
	},nil
}


func (ath* AttackTargetsHandler) findAttackTargetEntry(cfgId string) *AttackTargetEntry {

	for _,entry:= range ath.attackTargetsConfig.AttackTargets {

		if entry.Id == cfgId {
			return entry
		}
	}

	return nil
}

func (ath *AttackTargetsHandler) makeAttackTargetsContentFromTemplate(request *model.AddAttackTargetsRequest, templateFPath string, fname string) ([]byte, error) {

	contentBuffer := bytes.NewBuffer(make([]byte,0))

	tp := template.New(fname)

	tp = tp.Funcs(template.FuncMap{
		"tostring": func(arr []string) string {

			arrStr := make([]string,0)

			for _,a := range arr {

				arrStr = append(arrStr,fmt.Sprintf(`"%s"`,a))
			}

			return fmt.Sprintf("[%s]",strings.Join(arrStr,","))
		},
	})

	tp,err:= tp.ParseFiles(templateFPath)

	if err!=nil {
		errS := fmt.Sprintf("Parse Attack Target Template file:%s ,err:[%v]",templateFPath,err)
		log.Errorf(errS)
		return nil,fmt.Errorf(errS)
	}

	if err = tp.Execute(contentBuffer,request);err!=nil {
		errS := fmt.Sprintf("Parse Attack Target Template file:%s ,err:[%v]",templateFPath,err)
		log.Errorf(errS)
		return nil,fmt.Errorf(errS)
	}

	return contentBuffer.Bytes(),nil
}

func (ath *AttackTargetsHandler) makeAttackTargetsContent(entry *AttackTargetEntry,request *model.AddAttackTargetsRequest) ([]byte,error) {

	fpath := fmt.Sprintf("%s%s%s",ath.attackTargetsConfig.StorePath,ath.attackTargetsConfig.StoreRPath,entry.Fname)

	if !fileutils.FileIsExisted(fpath){

		errS := fmt.Sprintf("The Attack Targets file:%s not existed!",fpath)
		log.Errorf(errS)

		return nil,fmt.Errorf(errS)
	}

	if entry.IsTemplate {
		return ath.makeAttackTargetsContentFromTemplate(request,fpath,entry.Fname)
	}else {

		return ioutil.ReadFile(fpath)
	}
}

func (ath *AttackTargetsHandler) addAttackTargetsContext(ctx *attackTargetsContext) {

	ath.lock.Lock()
	defer ath.lock.Unlock()
	ctx.listNode = ath.requests.PushBack(ctx)

	ath.waitingRequests = ath.waitingRequests+1

}

func (ath *AttackTargetsHandler) removeAttackTargetsContext(ctx *attackTargetsContext) {

	ath.lock.Lock()
	defer ath.lock.Unlock()

	ath.requests.Remove(ctx.listNode)
	ath.waitingRequests = ath.waitingRequests-1

}

func (ath *AttackTargetsHandler) findAttackTargetsContext(nodeId string) *attackTargetsContext {

	ath.lock.Lock()
	defer ath.lock.Unlock()

	for listNode := ath.requests.Front();listNode!=nil;listNode = listNode.Next() {

		ctx := listNode.Value.(*attackTargetsContext)

		 if fetched,ok := ctx.nodesMap[nodeId];ok&&!fetched {

		 	return ctx
		 }
	}

	//not existed
	return nil
}

func (ath *AttackTargetsHandler) fetchedAttackTargetsContext(nodeId string,ctx *attackTargetsContext){

	ath.lock.Lock()
	defer ath.lock.Unlock()

	ctx.nodesMap[nodeId] = true

	ctx.waitedFetchs = ctx.waitedFetchs-1

	if ctx.waitedFetchs <= 0 {

		ctx.waitChan<-true

	}
}

func (ath *AttackTargetsHandler) makeAttackTargetsContext(request *model.AddAttackTargetsRequest,content []byte) *attackTargetsContext {

	attackTargets := &model.AttackTargets{
		Name:        fmt.Sprintf("%s_%d",request.Name,time.Now().UnixNano()),
		Size:        uint64(len(content)),
		AttackTypes: request.AttackTypes,
		Content:     content,
	}

	nodesMap := make(map[string]bool)
	for _,nodeId := range request.NodeIds {
		nodesMap[nodeId] = false
	}

	attackTargetsCtx := &attackTargetsContext{
		targets:      attackTargets,
		nodesMap:     nodesMap,
		waitedFetchs: len(request.NodeIds),
		waitChan:     make(chan bool),
	}

	ath.addAttackTargetsContext(attackTargetsCtx)

	return attackTargetsCtx
}

func (ath *AttackTargetsHandler) waitAttackTargetsFetched(requestCtx *attackTargetsContext) *model.AddAttackTargetsReply{

	var msg string
	tm := time.NewTicker(time.Duration(ath.waitFetchTimeout)*time.Millisecond)

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

		ath.removeAttackTargetsContext(requestCtx)

		return &model.AddAttackTargetsReply{
			Status:       0,
			Message:      msg,
			FetchedNodes: fetchedNodes,
		}
}

func (ath *AttackTargetsHandler) AddAttackTargetsHandle(request *model.AddAttackTargetsRequest) *model.AddAttackTargetsReply {

	log.Infof("Accept an add attack targets request:%v",request)

	if request.NodeIds == nil || len(request.NodeIds) == 0 {

		errs := fmt.Sprintf("Please Specify the nodes that accept this attack targets:%d",ath.waitingRequests)
		log.Errorf(errs)

		return &model.AddAttackTargetsReply{
			Status:       -1,
			Message:      errs,
			FetchedNodes: []string{},
		}
	}
	if ath.waitingRequests+1>ath.queueCapacity {

		errs := fmt.Sprintf("Too many add attack targets request to waiting,current requests:%d",ath.waitingRequests)
		log.Errorf(errs)

		return &model.AddAttackTargetsReply{
			Status:       -1,
			Message:      errs,
			FetchedNodes: []string{},
		}
	}

	entry := ath.findAttackTargetEntry(request.CfgId)

	if entry == nil {
		errS := fmt.Sprintf("cannot find attack targets config entry,cfgId:%s",request.CfgId)
		log.Errorf(errS)

		return &model.AddAttackTargetsReply{
			Status:       -1,
			Message:     errS ,
			FetchedNodes: []string{},
		}
	}

	content,err:= ath.makeAttackTargetsContent(entry,request)
	if err!=nil {

		log.Errorf("%v",err)
		return &model.AddAttackTargetsReply{
			Status:       -1,
			Message:      err.Error(),
			FetchedNodes: []string{},
		}
	}

	ctx := ath.makeAttackTargetsContext(request,content)

	//wait request to be accept by all nodes
	return ath.waitAttackTargetsFetched(ctx)

}

func (ath *AttackTargetsHandler) FetchAttackTargets(request *model.FetchAttackTargetsRequest) *model.AttackTargets {

	ctx := ath.findAttackTargetsContext(request.NodeId)

	if ctx == nil {

		//log.Infof("no find some attack targets to be added for this node:%s",request.NodeId)
		return &model.AttackTargets{
			Name:        "",
			Size:        0,
			AttackTypes: []string{},
			Content:     []byte{},
		}
	}

	log.Infof("find some attack targets to be added for this node:%s",request.NodeId)

	ath.fetchedAttackTargetsContext(request.NodeId,ctx)

	return ctx.targets
}