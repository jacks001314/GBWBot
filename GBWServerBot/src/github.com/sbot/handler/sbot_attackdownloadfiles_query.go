package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"strings"
)

func (sqh *SbotQueryHandler) queryStringOfAttackDownloadFiles(query *model.AttackedNodeDownloadFileQuery) string {

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


func (sqh *SbotQueryHandler) makeAttackedNodeDownloadFileMessage(entry *store.ResultEntry) (*model.AttackedNodeDownloadFileMessage,error) {

	var request AttackFileDownloadRequest

	if err := json.Unmarshal([]byte(entry.Value),&request);err!=nil {
		errS := fmt.Sprintf("Cannot json decode from json string:%s,failed:%v for attacked node download files",entry.Value,err)

		log.WithField("id",entry.Key).Errorln(errS)

		return nil,fmt.Errorf(errS)
	}

	return &model.AttackedNodeDownloadFileMessage{
		UserId: 	  request.UserId,
		TaskId:       request.TaskId,
		NodeId:       request.NodeId,
		Url:          request.Url,
		Fname:        request.Fname,
		AttackType:   request.AttackType,
		AttackIP:     request.AttackIP,
		TargetIP:     request.TargetIP,
		TargetPort:   int32(request.TargetPort),
		TargetOutIP:  request.TargetOutIP,
		DownloadTool: request.DownloadTool,
		UserAgent:    request.UserAgent,
		Time: request.Time,
	},nil
}

func (sqh *SbotQueryHandler) AttackedNodeDownloadFilesQueryHandle(query *model.AttackedNodeDownloadFileQuery) (*model.AttackedNodeDownloadFileReply,error) {

	db := sqh.attackedDownloadFilesDB

	queryString := sqh.queryStringOfAttackDownloadFiles(query)

	start,end := sqh.getTimeRange(query.Time.Start,query.Time.End)

	results,err := db.Query(queryString,[2]uint64{start,end},&store.Pageable{
		Page:  query.Page.Page,
		Size:  query.Page.Size,
		ISDec: true,
	})

	if err!=nil {
		errS := fmt.Sprintf("Query Attacked Node download files  database with query:%v failed:%v",query,err)
		log.Errorln(errS)
		return nil,fmt.Errorf(errS)
	}

	attackedNodeDownloaFilesReply := &model.AttackedNodeDownloadFileReply{
		Page: &model.PageMessage{
			Page:  results.Page,
			Size:  results.Size,
			TotalPage: results.TPage,
			TotalNum: results.TNum,
		},
		DownloadFiles: make([]*model.AttackedNodeDownloadFileMessage,0),
	}

	for _,entry := range results.Results {

		if adf,err:= sqh.makeAttackedNodeDownloadFileMessage(entry);err == nil {

			attackedNodeDownloaFilesReply.DownloadFiles = append(attackedNodeDownloaFilesReply.DownloadFiles,adf)
		}
	}

	return attackedNodeDownloaFilesReply,nil

}
