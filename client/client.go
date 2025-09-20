package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func Startclient() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error while connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to the server ")

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Println("Enter the msg:")
		var message string
		message,err = reader.ReadString('\n')
		if err != nil{
			fmt.Println("Error while reading string")
			continue
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error while sending message to server")
			continue
		}

	}
}
