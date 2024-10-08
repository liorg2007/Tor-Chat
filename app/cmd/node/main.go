package main

import (
	"bufio"
	"fmt"
	"log"
	"marshmello/pkg/session"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

const REDIS_ADDR = "addr"

var sm session.SessionManager

// Router function to redirect paths to their corresponding handlers
func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/get-aes":
		//GetAesHandler(w, r, sm)
	case "/set-redirect":
		//SetRedirectHandler(w, r, sm)
	case "/redirect":
		//RedirectHandler(w, r, sm)
	default:
		http.NotFound(w, r) // Default case for undefined paths
	}
}

// Function to handle incoming requests with goroutines
func handleRequest(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()
	router(w, r) // Redirect request to the appropriate handler
	fmt.Printf("Finished handling request from %s for %s\n", r.RemoteAddr, r.URL.Path)
}

// Function to listen for "EXIT" command and close the server
func consoleInput(listener *net.Listener, shutdown chan bool) {
	// Wait for EXIT command
	fmt.Println("Type 'EXIT' to close: ")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToUpper(text)) == "EXIT" {
			fmt.Println("Shutting down server...")
			close(shutdown)
			(*listener).Close() // Close the server listener
			return
		}
	}
}

func main() {
	// Create a WaitGroup to manage multiple goroutines
	var wg sync.WaitGroup
	sm = *session.NewSessionManager(REDIS_ADDR)

	// Create a channel to signal server shutdown
	shutdown := make(chan bool)

	// Create a new HTTP server and a handler function for all requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		go handleRequest(w, r, &wg) // Handle each request in a new goroutine
	})

	// Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

	// Start the HTTP server
	go func() {
		log.Println("Starting server on :8080")
		if err := http.Serve(listener, nil); err != nil {
			if err.Error() == "use of closed network connection" {
				log.Println("Server closed.")
			} else {
				log.Fatal("Server error: ", err)
			}
		}
	}()

	// Start listening for console input to shut down the server
	go consoleInput(&listener, shutdown)

	// Wait for the shutdown signal
	<-shutdown

	// Wait for all goroutines to finish before shutting down
	wg.Wait()
	log.Println("All requests have been processed. Server is now shut down.")
}
