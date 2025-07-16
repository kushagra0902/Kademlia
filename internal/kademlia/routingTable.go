package kademlia

import (
	"fmt"
	"math/big"
	"sort"
	"time"
)


const (
	IDlength = 160
	k        = 20
)

type Kbucket struct {
	Nodes []*KnownNodes
}

type RoutingTable struct {
	SelfID   [20]byte
	KBucktes [IDlength]*Kbucket
}

type KnownNodes struct {
	IP       string
	Port     string
	LastSeen time.Time
	NodeID   [20]byte
}

func XOR(idA [20]byte,idB [20]byte)[20]byte{
	var result [20]byte
	for i:=range 20{
		result[i]= idA[i]^idB[i]
	}
	return result
}

func Ownership(selfID [20]byte, targetID [20]byte) int{

// Each bucket in the routing table holds nodes within a certain XOR distance range from selfID.
// Buckets are indexed by the index of the first differing bit (from left to right).
// For example: if node A and node B differ at the 3rd bit, you put B in bucket 3 of A.

//This func compares a target node with SelfID and calculates bucket undex in which it should be stored.As every node has 160 K buckets in their own touting table and thus we have to calc relative index for each node.

//so in one sense from left we comapre first different Bit in both IDs and that represnet in which bucket it should go
var bucket_no int

xor := XOR(selfID, targetID)
	for i := 0; i < IDlength; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		if (xor[byteIndex]>>bitIndex)&1 == 1 {
			bucket_no=i
		}
	}
 // returns smallest index to be used 
return bucket_no
}

func NewRoutingTable(node Node)*RoutingTable{

	var rt RoutingTable
	rt.SelfID = node.ID
	for i := range IDlength{
		rt.KBucktes[i] = &Kbucket{}
	}
 return &rt
}

func (rt *RoutingTable) Update(nd Node) {
	node := KnownNodes{
		NodeID: nd.ID,
		IP: nd.IP,
		LastSeen: nd.LastSeen,
		Port: nd.Port,
	}
	bucket_idx := Ownership(rt.SelfID, node.NodeID)
	bucket := rt.KBucktes[bucket_idx]

	if bucket == nil {
		bucket = &Kbucket{}
		rt.KBucktes[bucket_idx] = bucket
	}

	//if alrready in the k-bucket, then it should just update its last seen
	for i, existing := range bucket.Nodes {
		if existing.NodeID == node.NodeID {
			bucket.Nodes = append(bucket.Nodes[:i], bucket.Nodes[i+1:]...)
			bucket.Nodes = append(bucket.Nodes, &node)
			return
		}
	}

	
	if len(bucket.Nodes) < k {
		bucket.Nodes = append(bucket.Nodes, &node)
		return
	}

	// If full, ping the least recently seen (first)
	target := Node{
		ID: bucket.Nodes[0].NodeID,
		IP:     bucket.Nodes[0].IP,
		Port:   bucket.Nodes[0].Port,
	}

	err := SendPing(rt.SelfID, target)
	if err != nil {
		bucket.Nodes = append(bucket.Nodes[:0], bucket.Nodes[1:]...) // Remove first
		bucket.Nodes = append(bucket.Nodes[1:], &node)
	}else{
		fmt.Println("K Bucket Full")
	}
	
	
}


func (rt *RoutingTable)FindClosest(key [20]byte) []*Node{
	var result []*Node
	var temp = make(map[*KnownNodes]big.Int)
	for _,buckets:=range rt.KBucktes{
		for _,nodes := range buckets.Nodes{
			dist := XOR(key, nodes.NodeID)
			dist_bigInt := new(big.Int).SetBytes(dist[:])
			temp[nodes] = *dist_bigInt
		}

	
	}
	type pair struct {
		Node     *KnownNodes
		Distance *big.Int
	}
	var pairs []pair
	for node, dist := range temp {
		d := new(big.Int).Set(&dist)
		pairs = append(pairs, pair{Node: node, Distance: d})
	}

	
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Distance.Cmp(pairs[j].Distance) == -1
	})

	// Pick top k closest nodes
	for i := 0; i < len(pairs) && i < k; i++ {
		node := &Node{
			ID:   pairs[i].Node.NodeID,
			IP:   pairs[i].Node.IP,
			Port: pairs[i].Node.Port,
		}
		result = append(result, node)
	}

	
	if len(result) <20{
		return result
	}

	return result
}

