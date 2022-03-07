package rservice

import (
	"context"
	"github.com/sbot/handler"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"io/ioutil"
)

type AttackPayloadService struct {

	service.UnimplementedAttackPayloadServiceServer
	handle *handler.AttackJarPayloadHandle

}

func NewAttackPayloadService(handle *handler.AttackJarPayloadHandle) *AttackPayloadService {

	return &AttackPayloadService{
		UnimplementedAttackPayloadServiceServer: service.UnimplementedAttackPayloadServiceServer{},
		handle:                                  handle,
	}
}

func (aps *AttackPayloadService) MakeJar(ctx context.Context, request *model.MakeJarAttackPayloadRequest) (*model.MakeJarAttackPayloadReply, error) {

	jarFile,err:= aps.handle.Handle(request)

	if err!=nil {

		return &model.MakeJarAttackPayloadReply{
			Status:  -1,
			Length:  0,
			Content: []byte{},
		},err
	}

	content,err:= ioutil.ReadFile(jarFile)

	if err!=nil {
		return &model.MakeJarAttackPayloadReply{
			Status:  -1,
			Length:  0,
			Content: []byte{},
		},err
	}

	return &model.MakeJarAttackPayloadReply{
		Status:  0,
		Length:  uint64(len(content)),
		Content: content,
	},nil
}