package network

import (
	"KademliaApply/internal/kademlia"
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

func SendPing(SelfID [20]byte, target kademlia.Node) error {
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
