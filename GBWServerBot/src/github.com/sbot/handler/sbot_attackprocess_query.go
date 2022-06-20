package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"strings"
)

func (sqh *SbotQueryHandler) queryStringOfAttackProcess(query *model.AttackProcessQuery) string {

	queryStringS := make([]string,0)

	if query.UserId!= "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["userId"]=="%s"`,query.UserId))
	}

	if query.TaskId!= "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["taskId"]=="%s"`,query.TaskId))
	}

	if query.NodeId != "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["nodeId"]=="%s"`,query.NodeId))
	}

	if query.AttackType != "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["attackType"]=="%s"`,query.AttackType))
	}


	if len(queryStringS) == 0 {

		return "1==1"
	}


	return strings.Join(queryStringS," and ")
}


func (sqh *SbotQueryHandler) makeAttackProcessMessage(entry *store.ResultEntry) (*model.AttackProcessMessage,error) {

	var attackProcess model.AttackProcessRequest

	if err := json.Unmarshal([]byte(entry.Value),&attackProcess);err!=nil {
		errS := fmt.Sprintf("Cannot json decode from json string:%s,failed:%v for attacke process",entry.Value,err)

		log.WithField("id",entry.Key).Errorln(errS)

		return nil,fmt.Errorf(errS)
	}

	return &model.AttackProcessMessage{
		UserId:     attackProcess.UserId,
		TaskId:     attackProcess.TaskId,
		NodeId:     attackProcess.NodeId,
		Time:       attackProcess.Time,
		TargetIP:   attackProcess.TargetIP,
		TargetHost: attackProcess.TargetHost,
		TargetPort: attackProcess.TargetPort,
		Proto:      attackProcess.Proto,
		App:        attackProcess.App,
		Os:         attackProcess.Os,
		AttackName: attackProcess.AttackName,
		AttackType: attackProcess.AttackType,
		Status:     attackProcess.Status,
		Payload:    attackProcess.Payload,
		Result:     attackProcess.Result,
		Details:    attackProcess.Details,
	},nil
}

func (sqh *SbotQueryHandler) AttackProcessQueryHandle(query *model.AttackProcessQuery) (*model.AttackProcessMessageReply,error) {

	db := sqh.attackProcessDB

	queryString := sqh.queryStringOfAttackProcess(query)

	start,end := sqh.getTimeRange(query.Time.Start,query.Time.End)

	results,err := db.Query(queryString,[2]uint64{start,end},&store.Pageable{
			Page:  query.Page.Page,
			Size:  query.Page.Size,
			ISDec: true,
	})

	if err!=nil {
		errS := fmt.Sprintf("Query Attacked Node  database with query:%v failed:%v",query,err)
		log.Errorln(errS)
		return nil,fmt.Errorf(errS)
	}

	attackProcessReply := &model.AttackProcessMessageReply{
		Page: &model.PageMessage{
				Page:  results.Page,
				Size:  results.Size,
				TotalPage: results.TPage,
				TotalNum: results.TNum,
			},

			Aps:  make([]*model.AttackProcessMessage,0),
	}

	for _,entry := range results.Results {

		if ap,err:= sqh.makeAttackProcessMessage(entry);err == nil {
			attackProcessReply.Aps = append(attackProcessReply.Aps,ap)
		}
	}

	return attackProcessReply,nil

}

