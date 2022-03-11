package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"strings"
	)

func (sqh *SbotQueryHandler) queryStringOfAttackedNode(query *model.AttackedNodeQuery) string {


	queryStringS := make([]string,0)

	if query.TaskId!= "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["taskId"]=="%s"`,query.TaskId))
	}

	if query.ParentNodeId != "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["pnodeId"]=="%s"`,query.ParentNodeId))
	}

	if query.Mac != "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["mac"]=="%s"`,query.Mac))
	}

	if query.AttackType != "" {

		queryStringS = append(queryStringS,fmt.Sprintf(`JsonValue["attackType"]=="%s"`,query.AttackType))
	}

	if len(queryStringS) == 0 {

		return "1==1"
	}


	return strings.Join(queryStringS," and ")
}


func (sqh *SbotQueryHandler) makeAttackedNodeMessage(entry *store.ResultEntry) (*model.AttackedNodeMessage,error) {

	var node model.CreateNodeRequest

	if err := json.Unmarshal([]byte(entry.Value),&node);err!=nil {
		errS := fmt.Sprintf("Cannot json decode from json string:%s,failed:%v for attacked node query",entry.Value,err)

		log.WithField("nodeId",entry.Key).Errorln(errS)

		return nil,fmt.Errorf(errS)
	}

	return &model.AttackedNodeMessage{
		TaskId:       node.TaskId,
		ParentNodeID: node.PnodeId,
		AttackType:    node.AttackType,
		NodeId:       entry.Key,
		Version:      node.Version,
		LocalIP:      node.LocalIP,
		OutIP:        node.OutIP,
		Mac:          node.Mac,
		Os:           node.Os,
		Arch:         node.Arch,
		User:         node.User,
		HostName:     node.HostName,
		Time:         node.Time,
		LastTime:     node.LastTime,
	},nil
}

func (sqh *SbotQueryHandler) AttackedNodeQueryHandle(query *model.AttackedNodeQuery) (*model.AttackedNodeReply,error) {

	db := sqh.attackedNodeDB

	if query.NodeId != "" {
		var node model.CreateNodeRequest

		if ok,err := db.Get(query.NodeId,&node);!ok {
			errS := fmt.Sprintf("Query Attacked Node Database by nodeId:%s,failed:%v ",query.NodeId,err)
			log.WithField("nodeId",query.NodeId).Errorln(errS)
			return nil,err
		}

		return &model.AttackedNodeReply{
			Page: &model.PageMessage{
				Page:  1,
				Size:  1,
				TotalNum: 1,
				TotalPage:1,
			},
			Nodes: []*model.AttackedNodeMessage{
				{
					TaskId:       node.TaskId,
					ParentNodeID: node.PnodeId,
					AttackType:   node.AttackType,
					NodeId:       query.NodeId,
					Version:      node.Version,
					LocalIP:      node.LocalIP,
					OutIP:        node.OutIP,
					Mac:          node.Mac,
					Os:           node.Os,
					Arch:         node.Arch,
					User:         node.User,
					HostName:     node.HostName,
					Time:         node.Time,
					LastTime:     node.LastTime,
				},
			},
		},nil

	}else {

		queryString := sqh.queryStringOfAttackedNode(query)
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

		attackedNodeReply := &model.AttackedNodeReply{
			Nodes: make([]*model.AttackedNodeMessage,0),
			Page:     &model.PageMessage{
				Page:  results.Page,
				Size:  results.Size,
				TotalPage: results.TPage,
				TotalNum: results.TNum,
			},
		}

		for _,entry := range results.Results {

			if attackedNode,err:= sqh.makeAttackedNodeMessage(entry);err == nil {

				attackedNodeReply.Nodes = append(attackedNodeReply.Nodes,attackedNode)
			}
		}

		return attackedNodeReply,nil
	}
}

