package kademlia

import (
	//"KademliaApply/internal/routing"

	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
)

const alpha = 3

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

	node.RoutingTable.Update(*SenderNode)

	return "ok"
}

func FindNodeHelper(entryNode *Node, targetID [20]byte)*Node{
	
	ShortList := entryNode.RoutingTable.FindClosest(targetID)
	//get top 20 closest nodes to target from its own routing table. If lesser than 20, then get the available amount.
	ShortList = append(ShortList, entryNode)
	ShortList = sortShortlistByDistance(ShortList, targetID) // sorted and length reduced to 20 or less (k)

	var queried  =  make (map[[20]byte]bool) 

	for {
		var parallel_req []*Node

		for _, node := range ShortList{
			if (!queried[node.ID] && len(parallel_req) < alpha){
				parallel_req = append(parallel_req, node)
				queried[node.ID] = true
			}
		}

		if len(parallel_req) == 0 {
			break
		}

		resultCh :=  make(chan []*Node, len(parallel_req))
		var wg  sync.WaitGroup
		for _, node := range parallel_req{
			wg.Add(1)
			go func(nd *Node){
				defer wg.Done()
				if (!queried[nd.ID]){
				fetched := nd.RoutingTable.FindClosest(targetID)
				resultCh <-fetched
				}
			}(node)
		}
		wg.Wait()
		close(resultCh)
		seen := make(map[[20]byte]bool)
		for _, node := range ShortList {
			seen[node.ID] = true
		}

		for fetched := range resultCh {
			for _, n := range fetched {
				if !seen[n.ID] {
					ShortList = append(ShortList, n)
					seen[n.ID] = true
				}
			}
		}

		// Re-sort and truncate to k
		ShortList = sortShortlistByDistance(ShortList, targetID)
		if len(ShortList) > k {
			ShortList = ShortList[:k]
		}

		// Optional: stop early if exact match found
		if ShortList[0].ID == targetID {
			return ShortList[0]
		}
	}

	// Return best effort result
	if len(ShortList) > 0 {
		return ShortList[0]
	}
	return nil


}

func sortShortlistByDistance(ShortList []*Node, targetID [20]byte) []*Node {
	type pair struct {
		node     *Node
		distance *big.Int
	}

	var pairs []pair
	for _, node := range ShortList {
		xor := XOR(targetID, node.ID)
		dist := new(big.Int).SetBytes(xor[:])
		pairs = append(pairs, pair{
			node:     node,
			distance: dist,
		})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].distance.Cmp(pairs[j].distance) == -1 // ascending order
	})

	var sorted []*Node
	for _, p := range pairs {
		sorted = append(sorted, p.node)
	}
	if len(sorted) < 20{
		return sorted
	}
	return sorted[:21]
}
