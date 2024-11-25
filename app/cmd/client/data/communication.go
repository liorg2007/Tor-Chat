package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendHttpRequest sends a POST request with JSON data to the given address and msgType path.
func SendHttpRequest(addr string, data interface{}, msgType string) ([]byte, error) {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	// Create the full URL by appending msgType to the address
	fullURL := fmt.Sprintf("http://%s/%s", addr, msgType)

	// Send the POST request
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code indicates an error
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error: %s", string(body))
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return respBody, nil
}

type MessageSender struct {
	Circuit list.List
}
