package servers

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var clients []net.Conn
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
		clientsMutex.Lock()
		clients = append(clients, conn)
		clientsMutex.Unlock()

		go handleclient(conn)

	}

		

}


func handleclient(conn net.Conn){
	defer func(){
		fmt.Println("Client disconnected",conn.RemoteAddr())
		removeclient(conn)
		conn.Close()
	}()

	buffer := make([]byte,1024)
	for{

		n,err := conn.Read(buffer)
		if err != nil{
			fmt.Println("Error at reading:",err)
			return 
		}
		message := string(buffer[:n])
		fmt.Printf("[%s]: %s\n", conn.RemoteAddr(), message)

		for _,cli := range clients {
			if cli != conn {
				_,err = cli.Write([]byte(fmt.Sprintf("[%s]: %s", conn.RemoteAddr(), message)))
				if err != nil{
					fmt.Println("Error while sendig to client",err)
				}
			}
		}

	}	


}



func removeclient(conn net.Conn){
	for i,c := range clients{
		if c == conn{
			clients = append(clients[:i],clients[i+1:]...)
			break
		} 
	}
}