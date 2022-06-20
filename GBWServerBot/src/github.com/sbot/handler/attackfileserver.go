package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/store"
	"github.com/sbot/utils/uuid"
	"time"
)

type AttackFileServerHandle struct {
	taskDb store.Store
	db store.Store
}

type AttackFileDownloadRequest struct {
	
	UserId string `json:"userId"`
	TaskId string `json:"taskId"`
	NodeId string `json:"nodeId"`
	AttackType string `json:"attackType"`

	Url   string `json:"url"`
	Fname string `json:"fname"`

	AttackIP   string `json:"attackIP"`
	TargetIP   string `json:"targetIP"`
	TargetPort int    `json:"targetPort"`

	TargetOutIP  string `json:"targetOutIP"`
	DownloadTool string `json:"downloadTool"`

	UserAgent string `json:"userAgent"`
	Time uint64 `json:"time"`
}

func NewAttackFileServerHandle(db store.Store,taskDB store.Store) *AttackFileServerHandle {

	return &AttackFileServerHandle{db: db,taskDb:taskDB}
}

func (afsh *AttackFileServerHandle) Handle(request *AttackFileDownloadRequest) error {

	key := fmt.Sprintf("AttackFile_%s", uuid.UUID())

	now := uint64(time.Now().UnixNano() / (1000 * 1000))

	request.Time = uint64(now)
	request.UserId = GetUserId(afsh.taskDb,request.TaskId)

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
