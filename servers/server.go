package servers

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/SanskarMali726/Broadcast-server/encryption"
	"github.com/joho/godotenv"
)

var clients = make(map[net.Conn]string)
var clientsMutex sync.Mutex

func Startserver() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error while loading env: %v", err)
	}

	PORT := os.Getenv("PORT")

	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("connect to the server on port", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error while accepting connection", err)
			return
		}
		fmt.Println("Client connected", conn.RemoteAddr())
		go handleclient(conn)
	}
}

func handleclient(conn net.Conn) {
	defer func() {
		fmt.Println("Client disconnected", conn.RemoteAddr())
		removeclient(conn)
		conn.Close()
	}()

	var username string
	for {
		conn.Write([]byte("Enter your username:"))
		
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading username:", err)
			return
		}
		
		name := strings.TrimSpace(string(buffer[:n]))
		if name == "" {
			conn.Write([]byte("Username cannot be empty. Try again:"))
			continue
		}
		
		if isNameTaken(name) {
			conn.Write([]byte("This name already taken.Try other"))
			continue
		} else {
			username = name
			clientsMutex.Lock()
			clients[conn] = username
			clientsMutex.Unlock()
			break
		}
	}
	
	var length uint32
	for {
	
		err := binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error while reading length", err)
			}
			return
		}

		buffer := make([]byte, length)
		totalRead := 0
		for totalRead < int(length) {
			n, err := conn.Read(buffer[totalRead:])
			if err != nil {
				fmt.Println("Error at reading:", err)
				return
			}
			totalRead += n
		}
		
		if len(buffer) < 12 {
			fmt.Println("Message too short")
			continue
		}
		
		key := buffer[:32]
		nonce := buffer[32:32+12]
		finalmsg := buffer[32+12:]

		msg, err := encryption.Decrypt(key, nonce, finalmsg)
		if err != nil {
			fmt.Println("Error while decrypting", err)
			return
		}

		fmt.Printf("[%s]: %s\n", username, string(buffer))

		clientsMutex.Lock()
		for co, cli := range clients {
			if cli != username {
				_, err = co.Write([]byte(fmt.Sprintf("[%s]: %s\n", username, string(msg))))
				if err != nil {
					fmt.Println("Error while sending to client", err)
				}
			}
		}
		clientsMutex.Unlock()
	}
}

func isNameTaken(name string) bool {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	
	for _, v := range clients {
		if v == name {
			return true
		}
	}
	return false
}

func removeclient(conn net.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, conn)
}