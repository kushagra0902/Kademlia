package kademlia

import (

	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

type Node struct {
	ID [20]byte 
	IP string
	Port string
	LastSeen time.Time
	RoutingTable RoutingTable
}

func GenerateNodeID(addr string) [20]byte {
	nodeID := sha1.Sum([]byte(addr))
	return nodeID
}

func (node *Node) PrintNodeID(){
	nodeIDstring:= hex.EncodeToString(node.ID[:])
	fmt.Println("[DEBUG]NodeID of the node is : ",nodeIDstring)
}


func NodeConstructor() *Node{
	addr := "123:456" // will get the IP port from stun server later
	var NewNode Node
	NewNode.ID = GenerateNodeID(addr)
	NewNode.PrintNodeID()
	// We will make a http client which will use the routes to listen to different rpcs. This is done in accordance to accomodate the other structure and to well define the files
	StartHttpServer(&NewNode)
	return &NewNode
}
