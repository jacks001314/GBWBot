package handler

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"github.com/sbot/utils/uuid"
)

type AttackTaskHandler struct {
	db store.Store
}

func (ath *AttackTaskHandler) doHandle(request model.CreateAttackTaskRequest) {

	taskId := fmt.Sprintf("taskId_%s", uuid.UUID())

}
