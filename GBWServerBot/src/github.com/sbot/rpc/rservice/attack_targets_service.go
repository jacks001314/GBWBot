package rservice

import (
	"context"
	"github.com/sbot/handler"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
)

type AttackTargetsService struct {

	service.UnimplementedAttackTargetsServiceServer

	handle *handler.AttackTargetsHandler

}


func NewAttackTargetsService(handle *handler.AttackTargetsHandler) *AttackTargetsService {

	return &AttackTargetsService{
		UnimplementedAttackTargetsServiceServer: service.UnimplementedAttackTargetsServiceServer{},
		handle:                                  handle,
	}
}

func (ats *AttackTargetsService) AddAttackTargets(_ context.Context, request *model.AddAttackTargetsRequest) (*model.AddAttackTargetsReply, error) {

	return ats.handle.AddAttackTargetsHandle(request),nil
}

func (ats *AttackTargetsService) FetchAttackTargets(_ context.Context, request *model.FetchAttackTargetsRequest) (*model.AttackTargets, error) {

	return ats.handle.FetchAttackTargets(request),nil
}
