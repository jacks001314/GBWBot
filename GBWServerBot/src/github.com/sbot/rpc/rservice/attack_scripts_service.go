package rservice

import (
	"context"
	"github.com/sbot/handler"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
)

type AttackScriptsService struct {

	service.UnimplementedAttackScriptsServiceServer

	handle *handler.AttackScriptsHandler

}

func NewAttackScriptsService(handle *handler.AttackScriptsHandler) *AttackScriptsService {

	return &AttackScriptsService{
		UnimplementedAttackScriptsServiceServer: service.UnimplementedAttackScriptsServiceServer{},
		handle:                                  handle,
	}
}

func (ass *AttackScriptsService) AddAttackScripts(_ context.Context, request *model.AddAttackScriptsRequest) (*model.AddAttackScriptsReply, error) {

	return ass.handle.AddAttackScriptsHandle(request),nil
}


func (ass *AttackScriptsService) FetchAttackScripts(request *model.FetchAttackScriptsRequest, assf service.AttackScriptsService_FetchAttackScriptsServer) error {

	scripts := ass.handle.FetchAttackScripts(request)

	if len(scripts) == 0 {

		//send an empty
		assf.Send(&model.AttackScripts{
			Name:         "",
			AttackType:   "",
			DefaultPort:  0,
			DefaultProto: "",
			Size:         0,
			Content:      []byte{},
			HasNext:      false,
		})

		return nil
	}

	for _,script:= range scripts {
		assf.Send(script)
	}

	return nil
}
