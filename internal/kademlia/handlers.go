package kademlia

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PingReq struct {
	From [20]byte // uneported fields dont get marshalled therefore made them public
	Req  string
	Addr string
}

type PingResponse struct {
	Status string `json:"status"`
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

func (node Node) ServeFind_Node() {

}

func (node Node) ServeFind_Key() {

}

func (node Node) ServeStore() {

}
