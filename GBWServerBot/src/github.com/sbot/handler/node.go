package handler

import "github.com/sbot/store"

type NodeHandler struct {
	dbnode          store.Store
	attackProcessDB store.Store
}

func NewNodeHandler(dbnode store.Store, attackProcessDb store.Store) *NodeHandler {

	return &NodeHandler{dbnode: dbnode,
		attackProcessDB: attackProcessDb,
	}
}
