package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"strings"
)


func (sqh *SbotQueryHandler) queryStringOfAttackTask(query *model.AttackTaskQuery) string {

	queryStringS := make([]string,0)

	if query.Name!= "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["name"]=="%s"`,query.Name))
	}

	if query.UserId != "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["userId"]=="%s"`,query.UserId))
	}

	if len(queryStringS) == 0 {

		return "1==1"
	}

	return strings.Join(queryStringS," and ")

}

func (sqh *SbotQueryHandler) getCbotPath(taskId string,osType model.OsType) string {

	return strings.ToLower(fmt.Sprintf("%s/cbot_%s.zip",taskId,model.OsType_name[int32(osType)]))
}

func (sqh *SbotQueryHandler) makeAttackTaskMessage(entry *store.ResultEntry) (*model.AttackTaskMessage,error) {

	var attackTask model.CreateAttackTaskRequest
	 if err := json.Unmarshal([]byte(entry.Value),&attackTask);err!=nil {
	 	errS := fmt.Sprintf("Cannot json decode from json string:%s,failed:%v for attack task query",entry.Value,err)

	 	log.WithField("taskId",entry.Key).Errorln(errS)

	 	return nil,fmt.Errorf(errS)
	 }

	return &model.AttackTaskMessage{
		TaskId:     entry.Key,
		Name:       attackTask.Name,
		UserId:     attackTask.UserId,
		Host:       attackTask.Host,
		Port:       attackTask.Port,
		User:       attackTask.User,
		Passwd:     attackTask.Passwd,
		PrivateKey: attackTask.PrivateKey,
		Os:         model.OsType_name[int32(attackTask.OsType)],
		Cbot:       attackTask.Cbot,
		CbotPath:   sqh.getCbotPath(entry.Key,attackTask.OsType),
	},nil
}

func (sqh *SbotQueryHandler) AttackTaskQueryHandle(query *model.AttackTaskQuery) (*model.AttackTaskReply,error) {

	db := sqh.attackTaskDB

	if query.TaskId != "" {
		var attackTask model.CreateAttackTaskRequest

		if ok,err := db.Get(query.TaskId,&attackTask);!ok {
			errS := fmt.Sprintf("Query Attack Task Database by taskId:%s,failed:%v ",query.TaskId,err)
			log.WithField("taskId",query.TaskId).Errorln(errS)
			return nil,err
		}

		return &model.AttackTaskReply{
			Messages: [] *model.AttackTaskMessage{
				{
				TaskId:     query.TaskId,
				Name:       attackTask.Name,
				UserId:     attackTask.UserId,
				Host:       attackTask.Host,
				Port:       attackTask.Port,
				User:       attackTask.User,
				Passwd:     attackTask.Passwd,
				PrivateKey: attackTask.PrivateKey,
				Os:         model.OsType_name[int32(attackTask.OsType)],
				Cbot:       attackTask.Cbot,
				CbotPath:   sqh.getCbotPath(query.TaskId,attackTask.OsType)},
			},
			Page: &model.PageMessage{
				Page:  1,
				Size:  1,
				TotalNum: 1,
				TotalPage:1,
			},
		},nil
	}else {

		queryString := sqh.queryStringOfAttackTask(query)
		start,end := sqh.getTimeRange(query.Time.Start,query.Time.End)

		results,err := db.Query(queryString,[2]uint64{start,end},&store.Pageable{
			Page:  query.Page.Page,
			Size:  query.Page.Size,
			ISDec: true,
		})

		if err!=nil {

			errS := fmt.Sprintf("Query Attack task database with query:%v failed:%v",query,err)
			log.Errorln(errS)
			return nil,fmt.Errorf(errS)

		}

		attackTaskReply := &model.AttackTaskReply{
			Messages: make([]*model.AttackTaskMessage,0),
			Page:     &model.PageMessage{
				Page:  results.Page,
				Size:  results.Size,
				TotalPage: results.TPage,
				TotalNum: results.TNum,
			},
		}

		for _,entry := range results.Results {

			if atm,err:= sqh.makeAttackTaskMessage(entry);err == nil {

				attackTaskReply.Messages = append(attackTaskReply.Messages,atm)
			}
		}

		return attackTaskReply,nil
	}
}
