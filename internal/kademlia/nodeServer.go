package kademlia

import (
	"fmt"
	"net/http"
)

func StartHttpServer(node *Node){
	
	http.HandleFunc("/ping", node.ServePing)
	err:=http.ListenAndServe(":5050", nil)
	if err!=nil{
		fmt.Println("[DEBUG]Error starting the server")
	}

}