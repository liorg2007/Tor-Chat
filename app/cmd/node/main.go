package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	// Middleware to log all requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse the form data and query parameters
			r.ParseForm()

			// Create the log message
			logMsg := fmt.Sprintf("\n\033[32m=== Request Details ===\n")
			logMsg += fmt.Sprintf("Method: %s\n", r.Method)
			logMsg += fmt.Sprintf("Path: %s\n", r.URL.Path)
			logMsg += fmt.Sprintf("Remote Address: %s\n", r.RemoteAddr)

			// Log headers
			logMsg += fmt.Sprintf("\nHeaders:\n")
			for key, values := range r.Header {
				logMsg += fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", "))
			}

			// Log query parameters
			if len(r.URL.Query()) > 0 {
				logMsg += fmt.Sprintf("\nQuery Parameters:\n")
				for key, values := range r.URL.Query() {
					logMsg += fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", "))
				}
			}

			// Log body if it exists
			if r.Body != nil && r.Header.Get("Content-Type") != "" {
				var bodyBytes []byte
				bodyBytes, _ = io.ReadAll(r.Body)
				// Restore the body for the actual handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				if len(bodyBytes) > 0 {
					logMsg += fmt.Sprintf("\nBody:\n  %s\n", string(bodyBytes))
				}
			}

			logMsg += "==================\033[0m\n"

			fmt.Print(logMsg)
			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/get-aes", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAesHandler(w, r, *sm)
	}).Methods("POST")

	r.HandleFunc("/set-redirect", func(w http.ResponseWriter, r *http.Request) {
		handlers.SetRedirectHandler(w, r, *sm)
	}).Methods("POST")

	r.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		handlers.RedirectHandler(w, r, *sm)
	}).Methods("POST")

	return r
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
