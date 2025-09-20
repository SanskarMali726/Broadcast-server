package main

import (
	"fmt"
	"os"

	"github.com/SanskarMali726/Broadcast-server/client"
	"github.com/SanskarMali726/Broadcast-server/servers"
)

func main(){
	
	

	if len(os.Args) < 2 {
		println("Usage: broadcast-server start|connect")
		return
	}

	command := os.Args[1]
	if command == "start" {
		servers.Startserver()
	}else if command == "connect"{
		client.Startclient()
	}else{
		fmt.Println("enter valid command")
	}
	
}

