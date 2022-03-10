package handler

import (
	"github.com/sbot/store"
	"time"
)

type SbotQueryHandler struct {

	attackTaskDB store.Store

	attackedNodeDB store.Store

	attackProcessDB store.Store

	attackedDownloadFilesDB store.Store


}


func NewSbotQueryHandler(attackTaskDB, attackedNodeDB, attackProcessDB, attackedDownloadFilesDB store.Store) *SbotQueryHandler {


	return &SbotQueryHandler{
		attackTaskDB:            attackTaskDB,
		attackedNodeDB:          attackedNodeDB,
		attackProcessDB:         attackProcessDB,
		attackedDownloadFilesDB: attackedDownloadFilesDB,
	}
}

func (sbh *SbotQueryHandler)getTimeRange(s,e uint64) (start, end uint64){

	var now uint64 = uint64(time.Now().UnixNano()/(1000*1000))

	if e == 0 || s>=e {

		return 0,now
	}

	return s,e
}




