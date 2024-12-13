package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode"
)

var (
	circuit   *MessageSender
	authToken string
)

func passwordChecker(password string) string {
	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return "Password is too short. It should be at least 8 characters long"
	} else if !hasUpperCase {
		return "Password should contain at least one uppercase letter"
	} else if !hasLowerCase {
		return "Password should contain at least one lowercase letter"
	} else if !hasDigit {
		return "Password should contain at least one digit"
	} else if !hasSpecial {
		return "Password should contain at least one special character"
	}

	return ""
}

func checkCredentials(req *AuthUserRequest) error {
	username := req.Username
	password := req.Password

	if len(username) < 5 || len(username) > 15 {
		return fmt.Errorf("{'detail' : 'Username must be 5-15 characters long'}")
	}

	for _, ch := range username {
		if !(unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '-') {
			return fmt.Errorf("{'detail' : 'Username can only have letters, digits, _ and -'}")
		}
	}

	errorMsg := passwordChecker(password)

	if errorMsg != "" {
		return fmt.Errorf("{'detail' : '%s", errorMsg)
	}

	return nil
}

func main() {
	circuitObj, err := CreateCircuit("http://localhost:8081/", "node2:8080", "node3:8080", "192.168.187.205:8080")
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

	err = checkCredentials(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	err = SendRegister(&circuit.Circuit, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
