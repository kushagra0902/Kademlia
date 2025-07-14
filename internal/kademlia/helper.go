package kademlia

import (
	"fmt"
	"time"
)

func PingHelper(node *Node, rt *RoutingTable, pingReq PingReq) string{
	fromNodeID := pingReq.From
	derivedNodeID := GenerateNodeID(pingReq.Addr)
	if fromNodeID != derivedNodeID {
		fmt.Println("NODE ID MISMATCH")
		return  "mismatch"
	}

	SenderNode := &Node{
		ID: fromNodeID,
		IP: pingReq.Addr,
		LastSeen: time.Now(),
	}

	_=SenderNode // later Will add func to add to routing table
	return "ok"
}