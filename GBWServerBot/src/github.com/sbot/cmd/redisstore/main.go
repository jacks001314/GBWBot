package main

import (
	"fmt"
	rstore "github.com/sbot/store"
	redisstore "github.com/sbot/store/redis"
	"github.com/sbot/utils/uuid"
	"time"
)

type Node struct {

	PNodeId string `json:"PNodeId"`
	NodeId string `json:"nodeId"`
	IP string `json:"ip"`
	Mac string `json:"mac"`
	
	Time uint64 `json:"time"`
}

func makeNode(ip string ) *Node {

	return &Node{
		PNodeId: uuid.UUID(),
		NodeId:  uuid.UUID(),
		IP:      ip,
		Mac:     "",
		Time:    uint64(time.Now().UnixNano()/(1000*1000)),
	}

}

func main(){

	var gnode Node

	var rdb redisstore.RedisStore

	cfg := &rstore.Config{
		DB:      "cbot",
		Table:   "node",
		Host:    "192.168.198.128",
		Port:    6379,
		User:    "",
		Pass:    "",
		Codes:   "",
		Timeout: 10000,
	}

	store,err:= rdb.Open(cfg)

	if err!=nil {
		fmt.Println(err)
		return
	}

	node := makeNode("192.168.1.151")
	store.Put(node.NodeId,node.Time,node)

	node = makeNode("192.168.1.152")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.152")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.152")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.153")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.154")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.154")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.155")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.155")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.155")
	store.Put(node.NodeId,node.Time,node)
	node = makeNode("192.168.1.155")
	store.Put(node.NodeId,node.Time,node)

	store.Get(node.NodeId,&gnode)

	r,err:= store.Query(`JsonValue["ip"]=="192.168.1.151"`,[2]uint64{0,node.Time},
		&rstore.Pageable{
			Page:  1,
			Size:  10,
			ISDec: true,
		})

	fmt.Println(r)

	qr,_:=store.Facet(`1==1`,`JsonValue["ip"]`,10,false)

	fmt.Println(qr)
}