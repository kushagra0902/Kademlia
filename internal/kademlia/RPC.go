package kademlia

import (
	
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PingReq struct {
	From [20]byte // uneported fields dont get marshalled therefore made them public
	Req  string
	Addr string
}

type PingResp struct {
	Status string `json:"status"`
}


func SendPing(SelfID [20]byte, target Node) error {
	address := "https://" + target.IP + ":" + target.Port + "/ping"
	var PingBody PingReq
	PingBody.From = SelfID
	PingBody.Req = "ping"

	jsonReq, err := json.Marshal(PingBody)

	if err != nil {
		fmt.Println("[DEBUG]Error marshalling the ping req json")
	}

	resp, err := http.Post(address, "application/json", bytes.NewBuffer(jsonReq))

	if err != nil {
		fmt.Println("[DEBUG]Error sending the ping req")
		return err
	}
	defer resp.Body.Close()
	var pingResp PingResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&pingResp)
	if err != nil {
		fmt.Println("[DEBUG]Error Decoding the ping response")
		return err
	}

	status := pingResp.Status
	if status != "ok" {
		err = fmt.Errorf("error pinging: status not ok")
		return err
	}

	fmt.Println("PING SUCCESS")
	return nil
}

func SendFindNode(SelfNode *Node, target *Node, targetID [20]byte) []*Node {
	req := FindNodeReq{
		Sender : SelfNode,
		TargetID: targetID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("[DEBUG] Failed to marshal FindNodeReq")
		return nil
	}

	url := fmt.Sprintf("http://%s:%s/find_node", target.IP, target.Port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[DEBUG] Failed to send FindNode to %s: %v\n", target.IP, err)
		return nil
	}
	defer resp.Body.Close()

	var nodeList []*Node
	err = json.NewDecoder(resp.Body).Decode(&nodeList)
	if err != nil {
		fmt.Println("[DEBUG] Failed to decode FindNode response")
		return nil
	}

	return nodeList
}


