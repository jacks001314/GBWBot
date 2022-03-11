package handler

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/utils/uuid"
	"time"
)

func (nh *NodeHandler) HandleCreateNode(request *model.CreateNodeRequest) (string, error) {

	now := uint64(time.Now().UnixNano() / (1000 * 1000))

	nodeId := fmt.Sprintf("node_%s", uuid.UUID())

	request.Time = now
	request.LastTime = now

	if err := nh.dbnode.Put(nodeId, now, request); err != nil {

		errS := fmt.Sprintf("Cannot write node:%s into database,err:%v", nodeId, err)

		log.Errorln(errS)

		return "", fmt.Errorf(errS)
	}

	log.Infof("Create nodeId:%s,parentNodeId:%s for request:%v\n", nodeId, request.PnodeId, request)

	return nodeId, nil
}
