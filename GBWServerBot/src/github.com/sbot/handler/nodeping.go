package handler

import (
	"github.com/sbot/proto/model"
	"time"
)

func (nh *NodeHandler) HandlePing(request *model.PingRequest) error {

	now := uint64(time.Now().UnixNano() / (1000 * 1000))

	var node model.CreateNodeRequest

	if ok, _ := nh.dbnode.Get(request.NodeId, &node); ok {

		node.LastTime = now

		//update
		nh.dbnode.Put(request.NodeId, now, &node)

	}

	log.WithTime(time.Now()).Infof("Receive Node:%s Ping", request.NodeId)

	return nil
}
