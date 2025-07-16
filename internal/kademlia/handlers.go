package kademlia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PingResponse struct {
	Status string `json:"status"`
}

type FindNodeReq struct {
	Sender *Node
	TargetID [20]byte `json:"targetid"`
}

func (node Node) ServePing(w http.ResponseWriter, r *http.Request) {
	var pingReq PingReq
	fmt.Println("Received Ping from:", r.RemoteAddr)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&pingReq)
	if err != nil {
		fmt.Println("[DEBUG]Error decoding the ping req at handler")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp := PingHelper(&node, &node.RoutingTable, pingReq)

	if resp != "ok" {
		http.Error(w, "NodeID Mismatch", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PingResponse{Status: "ok"})
}

func (node *Node) ServeFind_Node(w http.ResponseWriter, r *http.Request) {

	var FindReq FindNodeReq
	fmt.Println("Received FindNode req from :", r.RemoteAddr)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&FindReq)
	if err != nil {
		fmt.Println("[DEBUG]Error decoding the ping req at handler")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	targetID := FindReq.TargetID
	var zeroID [20]byte // imp to do like this as in this case len of targetID is always 20 even if it is empty as we already fixed its length
	if targetID == zeroID{
		http.Error(w, "No Target ID ", http.StatusUnauthorized)
	}
	sender := Node{
    ID: FindReq.Sender.ID,
    IP: FindReq.Sender.IP,
    Port: FindReq.Sender.Port,
    LastSeen: time.Now(),
}
	node.RoutingTable.Update(sender)
	requiredNode := FindNodeHelper(node, targetID)

	json.NewEncoder(w).Encode(requiredNode)
}

	

func (node Node) ServeFind_Key() {

}

func (node Node) ServeStore() {

}
