package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/utils/uuid"
	"time"
)

func (nh *NodeHandler) HandleAttackProcess(request *model.AttackProcessRequest) error {

	key := fmt.Sprintf("AttackProcess_%s", uuid.UUID())

	now := uint64(time.Now().UnixNano() / (1000 * 1000))

	jdata, _ := json.Marshal(request)

	if err := nh.attackProcessDB.Put(key, now, request); err != nil {

		errS := fmt.Sprintf("Cannot write attack process into database,failed:%v", err)

		log.Error(errS)

		return fmt.Errorf(errS)
	}

	log.WithTime(time.Now()).
		WithField("taskId", request.TaskId).
		WithField("nodeId", request.NodeId).
		Infof("Receive a attack process:%s", string(jdata))

	return nil
}
