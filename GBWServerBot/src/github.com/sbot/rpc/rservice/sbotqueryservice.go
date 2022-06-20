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


func (sqs *SbotQueryService) FacetAttackTasks(_ context.Context, request *model.FacetRequest) (*model.FacetReply, error) {

	return sqs.handle.FacetHandle(sqs.handle.GetAttackTasksDB(),request)
}

func (sqs *SbotQueryService) CountAttackTasks(_ context.Context, request *model.CountRequest) (*model.Count, error) {
	return sqs.handle.CountHandle(sqs.handle.GetAttackTasksDB(),request)
}

func (sqs *SbotQueryService) QueryAttackedNodes(_ context.Context, query *model.AttackedNodeQuery) (*model.AttackedNodeReply, error) {

	return sqs.handle.AttackedNodeQueryHandle(query)

}

func (sqs *SbotQueryService) FacetAttackedNodes(_ context.Context, request *model.FacetRequest) (*model.FacetReply, error) {
	return sqs.handle.FacetHandle(sqs.handle.GetAttackedNodesDB(),request)
}

func (sqs *SbotQueryService) CountAttackedNodes(_ context.Context, request *model.CountRequest) (*model.Count, error) {
	return sqs.handle.CountHandle(sqs.handle.GetAttackedNodesDB(),request)
}

func (sqs *SbotQueryService) QueryAttackProcess(_ context.Context, query *model.AttackProcessQuery) (*model.AttackProcessMessageReply, error) {

	return sqs.handle.AttackProcessQueryHandle(query)
}

func (sqs *SbotQueryService) FacetAttackProcess(_ context.Context, request *model.FacetRequest) (*model.FacetReply, error) {
	return sqs.handle.FacetHandle(sqs.handle.GetAttackProcessDB(),request)
}

func (sqs *SbotQueryService) CountAttackProcess(_ context.Context, request *model.CountRequest) (*model.Count, error) {
	return sqs.handle.CountHandle(sqs.handle.GetAttackProcessDB(), request)
}

func (sqs *SbotQueryService) QueryAttackedDownloadFiles(_ context.Context, query *model.AttackedNodeDownloadFileQuery) (*model.AttackedNodeDownloadFileReply, error) {

	return sqs.handle.AttackedNodeDownloadFilesQueryHandle(query)

}

func (sqs *SbotQueryService) FacetAttackedDownloadFiles(_ context.Context, request *model.FacetRequest) (*model.FacetReply, error) {
	return sqs.handle.FacetHandle(sqs.handle.GetAttackedDownloadFileDB(),request)
}

func (sqs *SbotQueryService) CountAttackedDownloadFiles(_ context.Context, request *model.CountRequest) (*model.Count, error) {
	return sqs.handle.CountHandle(sqs.handle.GetAttackedDownloadFileDB(),request)

}