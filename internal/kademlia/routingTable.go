package kademlia

import "time"

const (
	IDlength   = 160
	k = 20
)

type Kbucket struct {
	Nodes [k]*KnownNodes
}

type RoutingTable struct {
	SelfID   [20]byte
	KBucktes [IDlength]*Kbucket
}

type KnownNodes struct {
	IP       string
	Port     string
	LastSeen time.Time
	NodeID [20]byte
}

func XOR(idA [20]byte, idB[20]byte)[20]byte{
	var result [20]byte
	for i:=0;i<20;i++{
		result[i]= idA[i]^idB[i]
	}
	return result
}

