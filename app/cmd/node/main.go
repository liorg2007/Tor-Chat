package main

import (
	"bufio"
	"fmt"
	"marshmello/pkg/networking"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	listener, err := networking.CreateListening(networking.CONN_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer listener.Close()

	var wg sync.WaitGroup
	shutdown := make(chan bool)

	wg.Add(1)
	go clientAccepter(&listener, shutdown, &wg)
	go consoleInput(&listener, shutdown)
	// Wait for all client handlers to finish
	wg.Wait()
	fmt.Println("Server shutdown complete")
}

func clientAccepter(listener *net.Listener, shutdown chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		conn, err := (*listener).Accept()
		if err != nil {
			select {
			case <-shutdown:
				return // Normal shutdown
			default:
				fmt.Println("Error accepting connection:", err)
			}
			continue
		}

		wg.Add(1)
		//go handleClient(conn, &wg)
		var b []byte
		conn.Read(b)
	}
}

func consoleInput(listener *net.Listener, shutdown chan bool) {
	// Wait for EXIT command
	fmt.Println("Type 'EXIT' to close: ")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToUpper(text)) == "EXIT" {
			fmt.Println("Shutting down server...")
			close(shutdown)
			(*listener).Close()
			return
		}
	}
}
