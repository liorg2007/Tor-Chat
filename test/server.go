package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :1234...")

	conn, err := listener.Accept()

	defer conn.Close()

	fmt.Println("A client connected!!")
}
