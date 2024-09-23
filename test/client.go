package client

import (
	"fmt"
	"net"
)

var addr string = "localhost"

func main() {
	var port string
	fmt.Println("Hello world, Enter the port: ")
	fmt.Scan(&port)

	fullAddr := net.JoinHostPort(addr, port)

	conn, err := net.Dial("tcp", fullAddr)

	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected!!!")
}
