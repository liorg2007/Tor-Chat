package networking

import (
	"encoding/base64"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// Test case for successfully sending and receiving base64-encoded data
func TestSendAndReceiveData(t *testing.T) {
	// Create a mock network connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Sample base64 data
	data := base64.StdEncoding.EncodeToString([]byte("Test data"))

	// Send data in a separate goroutine
	go func() {
		err := SendData(client, data)
		if err != nil {
			t.Errorf("Failed to send data: %v", err)
		}
	}()

	// Read the data on the server side
	receivedData, err := ReceiveData(server)
	if err != nil {
		t.Errorf("Failed to receive data: %v", err)
	}

	// Verify the received data matches the sent data
	if receivedData != data {
		t.Errorf("Expected %v, got %v", data, receivedData)
	}
}

// Test case for sending and receiving an empty string
func TestSendAndReceiveEmptyData(t *testing.T) {
	// Create a mock network connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Empty data
	data := ""

	// Send data in a separate goroutine
	go func() {
		err := SendData(client, data)
		if err != nil {
			t.Errorf("Failed to send data: %v", err)
		}
	}()

	// Read the data on the server side
	receivedData, err := ReceiveData(server)
	if err != nil {
		t.Errorf("Failed to receive data: %v", err)
	}

	// Verify the received data matches the empty string
	if receivedData != data {
		t.Errorf("Expected empty string, got %v", receivedData)
	}
}

// Test case for handling invalid base64 data
func TestReceiveDataInvalidBase64(t *testing.T) {
	// Create a mock network connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Send invalid (non-base64) data
	go func() {
		length := uint32(len("Invalid base64 data"))

		// Create a byte slice for the length prefix (4 bytes for uint32 in big endian)
		lengthPrefix := make([]byte, 4)

		binary.BigEndian.PutUint32(lengthPrefix, length)
		client.Write(lengthPrefix)
		client.Write([]byte("Invalid base64 data"))
	}()

	// Attempt to receive the invalid data
	_, err := ReceiveData(server)
	if err == nil {
		t.Errorf("Expected error for invalid base64 data, but got none")
	}
}

// Test case for handling a network timeout
func TestSendDataWithTimeout(t *testing.T) {
	// Create a mock network connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Set a short timeout for the client-side write operation
	client.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))

	// Try to send data that will be delayed (simulate network issues)
	data := base64.StdEncoding.EncodeToString([]byte("Delayed data"))
	err := SendData(client, data)

	// Expect an error due to timeout
	if err == nil {
		t.Errorf("Expected a timeout error, but got none")
	}
}

// Test case for large data transfer
func TestSendAndReceiveLargeData(t *testing.T) {
	// Create a mock network connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Generate a large string of repeated characters and base64-encode it
	largeData := make([]byte, 10000) // 10KB of data
	for i := range largeData {
		largeData[i] = 'A'
	}
	data := base64.StdEncoding.EncodeToString(largeData)

	// Send large data in a separate goroutine
	go func() {
		err := SendData(client, data)
		if err != nil {
			t.Errorf("Failed to send data: %v", err)
		}
	}()

	// Receive the large data on the server side
	receivedData, err := ReceiveData(server)
	if err != nil {
		t.Errorf("Failed to receive data: %v", err)
	}

	// Verify the received data matches the sent data
	if receivedData != data {
		t.Errorf("Expected large data, but received different data")
	}
}

// Test case for ensuring proper closure of connections after sending data
func TestSendDataConnectionClosed(t *testing.T) {
	// Create a mock network connection using a pipe
	server, client := net.Pipe()

	// Sample base64 data
	data := base64.StdEncoding.EncodeToString([]byte("Test closing connection"))

	// Send data in a separate goroutine
	go func() {
		err := SendData(client, data)
		if err != nil {
			t.Errorf("Failed to send data: %v", err)
		}
		client.Close() // Close the client connection
	}()

	// Read the data on the server side
	receivedData, err := ReceiveData(server)
	if err != nil {
		t.Errorf("Failed to receive data: %v", err)
	}

	// Verify if the received data matches the sent data
	if receivedData != data {
		t.Errorf("Expected %v, got %v", data, receivedData)
	}

	// Verify if the connection is closed on the server side
	_, err = ReceiveData(server)
	if err == nil {
		t.Errorf("Expected error for receiving data after connection closure, but got none")
	}
}
