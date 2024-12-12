package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	circuit   *MessageSender
	authToken string
)

func main() {
	circuitObj, err := CreateCircuit("http://localhost:8081/", "node2:8080", "node3:8080", "10.10.246.33:8080")
	if err != nil {
		log.Fatal("Cant connect")
	}

	circuit = &circuitObj

	// Set up routes
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/send-message", sendMessageHandler)
	http.HandleFunc("/receive-messages", receiveMessagesHandler)

	// Start the server
	fmt.Println("Server starting on :1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = SendRegister(&circuit.Circuit, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := SendLogin(&circuit.Circuit, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Login failed: %v", err), http.StatusUnauthorized)
		return
	}

	// Store the token globally (in a real application, you'd use a more secure method)
	authToken = token

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is logged in
	if authToken == "" {
		http.Error(w, "Unauthorized: Please login first", http.StatusUnauthorized)
		return
	}

	var req SendMessageStruct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use the stored token
	req.Token = authToken

	err = SendMessage(&circuit.Circuit, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

func receiveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is logged in
	if authToken == "" {
		http.Error(w, "Unauthorized: Please login first", http.StatusUnauthorized)
		return
	}

	messages, err := ReceiveMessages(&circuit.Circuit, GetMessagees{Token: authToken})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to receive messages: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
