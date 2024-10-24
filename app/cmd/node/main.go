package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"marshmello/pkg/handlers"
	"marshmello/pkg/session"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var sm *session.SessionManager = nil

// WriteErrorResponse writes a standard JSON error response to the http.ResponseWriter.
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Router function to redirect paths to their corresponding handlers
func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/get-aes", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAesHandler(w, r, *sm)
	}).Methods("GET")

	r.HandleFunc("/set-redirect", func(w http.ResponseWriter, r *http.Request) {
		handlers.SetRedirectHandler(w, r, *sm)
	}).Methods("GET")

	r.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		handlers.RedirectHandler(w, r, *sm)
	}).Methods("GET")

	return r
}

// Function to handle incoming requests with goroutines
func handleRequest(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()
	router().ServeHTTP(w, r) // Use Gorilla Mux to handle the request
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
	var err error

	// Connecting to Redis service
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "localhost" // default fallback
	}
	if redisPort == "" {
		redisPort = "6379" // default fallback
	}

	// Uncomment the following lines to initialize the session manager

	sm, err = session.NewSessionManager(fmt.Sprintf("%s:%s", redisHost, redisPort))
	if err != nil {
		log.Fatal("Error connecting to Redis service: ", err)
		return
	}

	log.Printf("Connected to Redis service on %s:%s", redisHost, redisPort)

	// Create a channel to signal server shutdown
	shutdown := make(chan bool)

	// Create a new Gorilla Mux router
	r := router()

	// Start listening for requests
	http.Handle("/", r)

	// Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting the server: ", err)
		return
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
