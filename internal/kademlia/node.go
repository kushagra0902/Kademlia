package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Node struct {
	ID [20]byte 
	IP string
	LastSeen int64
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

func (node *Node) ServeRPC(rpc string){
	switch rpc {
	case "ping" :
		fmt.Println("[DEBUG]Serving ping RPC")
		node.ServePing()
	case "find_node" :
		fmt.Println("[DEBUG]Serving find_node RPC")
		node.ServeFind_Node()
	case "Store":
		fmt.Println("[DEBUG]Serving Store RPC")
		node.ServeStore()

	case "find_key":
		fmt.Println("[DEBUG]Serving find_key RPC")
		node.ServeFind_Key()
	}

}

func (node Node)ServePing(){

}

func (node Node)ServeFind_Node(){

}

func(node Node)ServeFind_Key(){

}

func (node Node)ServeStore(){

}

func NodeConstructor() *Node{
	addr := "123:456" // will get the IP port from stun server later
	var NewNode Node
	NewNode.ID = GenerateNodeID(addr)
	NewNode.PrintNodeID()
	rpc := "ping" // hardcoded as of now, later fetch from network
	NewNode.ServeRPC(rpc)

	return &NewNode
}
