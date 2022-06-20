package handler

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/store"
	"sync"
)

var queryLock = &sync.Mutex{}
var taskId2UserId = make(map[string]string)

/*Get userId from attack task table*/
func GetUserId(db store.Store,taskId string) string {

	var attackTask model.CreateAttackTaskRequest

	queryLock.Lock()
	defer queryLock.Unlock()

	if uid,ok := taskId2UserId[taskId];ok {
		return uid
	}

	if ok,err := db.Get(taskId,&attackTask);!ok {
		errS := fmt.Sprintf("Query Attack Task Database by taskId:%s,failed:%v ",taskId,err)
		log.WithField("taskId",taskId).Errorln(errS)
		return "nil"
	}

	taskId2UserId[taskId] = attackTask.UserId

	return attackTask.UserId
}
