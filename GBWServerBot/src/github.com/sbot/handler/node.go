package handler

import "github.com/sbot/store"

type NodeHandler struct {
	dbnode          store.Store
	taskDB			store.Store
	attackProcessDB store.Store
}

func NewNodeHandler(dbnode store.Store, attackProcessDb store.Store,taskDB store.Store) *NodeHandler {

	return &NodeHandler{dbnode: dbnode,
		attackProcessDB: attackProcessDb,
		taskDB:taskDB,
	}
}
