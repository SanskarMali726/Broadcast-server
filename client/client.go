package client

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/SanskarMali726/Broadcast-server/encryption"
)

func Startclient() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error while connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to the server")
	
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading username prompt:", err)
		return
	}
	fmt.Print(string(buffer[:n]))


	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	_, err = conn.Write([]byte(username))
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error while reading msg from:", conn.RemoteAddr(), "Error is:", err)
				}
				return
			}
			message := string(buffer[:n])
			fmt.Print(message)
		}
	}()

	
	for {
		var message string
		message, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading string")
			continue
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		key := make([]byte,32)
		_,err = rand.Read(key)
		if err != nil{
			fmt.Println(err)
			return 
		}

		msg, nonce, err := encryption.Encrypt(key, message)
		if err != nil {
			fmt.Println(err)
			return
		}

		finalmsg := append(append(key,nonce...),msg...)
		length := uint32(len(finalmsg))


		err = binary.Write(conn, binary.BigEndian, length)
		if err != nil {
			fmt.Println("Error while writing length", err)
			return
		}

		_, err = conn.Write(finalmsg)
		if err != nil {
			fmt.Println("Error while sending message to server")
			continue
		}
	}
}