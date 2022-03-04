package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/server"
	"github.com/sbot/store"
	"github.com/sbot/utils/uuid"
	"time"
)

type AttackFileServerHandle struct {
	db store.Store
}

func NewAttackFileServerHandle(db store.Store) *AttackFileServerHandle {

	return &AttackFileServerHandle{db: db}
}

func (afsh *AttackFileServerHandle) Handle(request *server.AttackFileDownloadRequest) error {

	key := fmt.Sprintf("AttackFile_%s", uuid.UUID())

	now := uint64(time.Now().UnixNano() / (1000 * 1000))

	jdata, _ := json.Marshal(request)

	if err := afsh.db.Put(key, now, request); err != nil {

		errS := fmt.Sprintf("Cannot write attack file download request into database,failed:%v", err)

		log.Error(errS)

		return fmt.Errorf(errS)
	}

	log.WithTime(time.Now()).
		WithField("taskId", request.TaskId).
		WithField("nodeId", request.NodeId).
		Infof("Receive and handle a attack file download request:%s", string(jdata))

	return nil
}
