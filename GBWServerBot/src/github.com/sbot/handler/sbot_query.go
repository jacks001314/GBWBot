package handler

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"time"
)

type SbotQueryHandler struct {

	attackTaskDB store.Store

	attackedNodeDB store.Store

	attackProcessDB store.Store

	attackedDownloadFilesDB store.Store

	sbot string
}


func NewSbotQueryHandler(sbot string,attackTaskDB, attackedNodeDB, attackProcessDB, attackedDownloadFilesDB store.Store) *SbotQueryHandler {


	return &SbotQueryHandler{
		attackTaskDB:            attackTaskDB,
		attackedNodeDB:          attackedNodeDB,
		attackProcessDB:         attackProcessDB,
		attackedDownloadFilesDB: attackedDownloadFilesDB,
		sbot: sbot,
	}
}

func (sbh *SbotQueryHandler) GetAttackTasksDB() store.Store{

	return sbh.attackTaskDB
}

func (sbh *SbotQueryHandler) GetAttackedNodesDB() store.Store{

	return sbh.attackedNodeDB
}

func (sbh *SbotQueryHandler) GetAttackProcessDB() store.Store{

	return sbh.attackProcessDB
}

func (sbh *SbotQueryHandler) GetAttackedDownloadFileDB() store.Store{

	return sbh.attackedDownloadFilesDB
}


func (sbh *SbotQueryHandler)getTimeRange(s,e uint64) (start, end uint64){

	var now uint64 = uint64(time.Now().UnixNano()/(1000*1000))

	if e == 0 || s>=e {

		return 0,now
	}

	return s,e
}

func (sqh *SbotQueryHandler) FacetHandle(db store.Store ,request *model.FacetRequest) (*model.FacetReply,error) {

	q := "1==1"

	results := &model.FacetReply{Items:make([]*model.FacetItem,0)}

	res,err := db.Facet(q,fmt.Sprintf(`JsonValue["%s"]`,request.Term), uint64(request.TopN),request.IsDec)

	if err!=nil {
		errS := fmt.Sprintf("Facet database with request:%v failed:%v",request,err)
		log.Errorln(errS)
		return nil,fmt.Errorf(errS)
	}

	for _,item := range res {

		results.Items = append(results.Items,&model.FacetItem{
			Key:   item.Key,
			Count: item.Count,
		})
	}

	return results,nil
}

func (sqh *SbotQueryHandler) CountHandle(db store.Store,request *model.Empty) (*model.Count,error) {

	return &model.Count{C:db.Count()},nil
}



