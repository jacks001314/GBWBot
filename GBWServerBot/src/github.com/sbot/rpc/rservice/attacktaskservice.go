package rservice

import (
	"context"
	"fmt"
	"github.com/sbot/handler"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
)

type AttackTaskService struct {
	service.UnimplementedAttackTaskServiceServer

	handle *handler.AttackTaskHandler
}

func NewAttackTaskService(handle *handler.AttackTaskHandler) *AttackTaskService {

	return &AttackTaskService{
		UnimplementedAttackTaskServiceServer: service.UnimplementedAttackTaskServiceServer{},
		handle:                               handle,
	}

}

func (ats *AttackTaskService) CreateAttackTask(ctx context.Context, request *model.CreateAttackTaskRequest) (*model.CreateAttackTaskReply, error) {

	log.Infof("Accept a create attack task request:%s", request)

	taskId, err := ats.handle.Handle(request)

	if err != nil {

		errS := fmt.Sprintf("Create attack task:%s failed:%v", request.Name, err)

		log.Error(errS)

		return &model.CreateAttackTaskReply{
			Status: -1,
			TaskId: "",
		}, fmt.Errorf(errS)
	}

	return &model.CreateAttackTaskReply{
		Status: 0,
		TaskId: taskId,
	}, nil

}
