package servers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var clients = make(map[net.Conn]string)
var clientsMutex sync.Mutex

func Startserver(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatalf("Error while loading env: %v",err)
	}

	PORT:=os.Getenv("PORT")

	l,err := net.Listen("tcp",PORT)
	if err != nil{
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("connect to the server on port ",PORT)

	for{

		conn,err := l.Accept()
		if err != nil{
			fmt.Println("Error while accepting connection",err)
			return
		}
		fmt.Println("Client connected",conn.RemoteAddr())
		go handleclient(conn)

	}

		

}


func handleclient(conn net.Conn){
	defer func(){
		fmt.Println("Client disconnected",conn.RemoteAddr())
		removeclient(conn)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	var username string
	for{
		conn.Write([]byte("Enter your username:"))
		name,_ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if isNameTaken(name){
			conn.Write([]byte("This name already taken.Try other"))
			continue
		}else{
			username = name
			clientsMutex.Lock()
			clients[conn] = username
			clientsMutex.Unlock()
			break
		}
	}

	buffer := make([]byte,1024)
	for{

		n,err := conn.Read(buffer)
		if err != nil{
			fmt.Println("Error at reading:",err)
			return 
		}
		
		message := string(buffer[:n])
		fmt.Printf("[%s]: %s\n",username, message)

		for co,cli := range clients {
			if cli != username {
				_,err = co.Write([]byte(fmt.Sprintf("[%s]: %s",username, message)))
				if err != nil{
					fmt.Println("Error while sendig to client",err)
				}
			}
		}

	}	

}

func isNameTaken(name string)bool{
	for _,v := range clients{
		if v == name{
			return true
		}
	}
	return false
}


func removeclient(conn net.Conn){
	delete(clients,conn)
}