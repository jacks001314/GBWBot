package rservice

import (
	"context"
	"github.com/sbot/handler"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
)

type SbotQueryService struct {

	service.UnimplementedSbotServiceServer

	handle *handler.SbotQueryHandler

}

func NewSbotQueryservice(handle *handler.SbotQueryHandler) *SbotQueryService {

	return &SbotQueryService{
		UnimplementedSbotServiceServer: service.UnimplementedSbotServiceServer{},
		handle:                         handle,
	}
}

func (sqs *SbotQueryService) QueryAttackTasks(_ context.Context, query *model.AttackTaskQuery) (*model.AttackTaskReply, error) {

	return sqs.handle.AttackTaskQueryHandle(query)
}

func (sqs *SbotQueryService) QueryAttackedNodes(_ context.Context, query *model.AttackedNodeQuery) (*model.AttackedNodeReply, error) {

	return sqs.handle.AttackedNodeQueryHandle(query)

}

func (sqs *SbotQueryService) QueryAttackProcess(_ context.Context, query *model.AttackProcessQuery) (*model.AttackProcessMessageReply, error) {

	return sqs.handle.AttackProcessQueryHandle(query)
}

func (sqs *SbotQueryService) QueryAttackedDownloadFiles(_ context.Context, query *model.AttackedNodeDownloadFileQuery) (*model.AttackedNodeDownloadFileReply, error) {

	return sqs.handle.AttackedNodeDownloadFilesQueryHandle(query)

}