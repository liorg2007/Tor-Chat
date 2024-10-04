package torpacket_test

import (
	"encoding/base64"
	"encoding/json"
	"marshmello/pkg/torpacket"
	"testing"
)

func TestSerializeGetAES(t *testing.T) {
	// Test correct serialization
	aesKey := make([]byte, torpacket.AES_KEY_LENGTH)
	for i := range aesKey {
		aesKey[i] = byte(i)
	}

	rawMsg, err := torpacket.SerializeGetAES(aesKey)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check the message code
	if rawMsg.Code != torpacket.GetAES {
		t.Errorf("expected message code %d, got %d", torpacket.GetAES, rawMsg.Code)
	}

	// Check the message content
	var msg torpacket.GetAesMsg
	err = json.Unmarshal([]byte(rawMsg.JsonData), &msg)
	if err != nil {
		t.Errorf("expected no error unmarshalling JSON, got %v", err)
	}

	expectedData := base64.StdEncoding.EncodeToString(aesKey)
	if msg.AesKey != expectedData {
		t.Errorf("expected aesKey %s, got %s", expectedData, msg.AesKey)
	}

	// Test incorrect AES key length
	shortKey := make([]byte, 10) // invalid length
	_, err = torpacket.SerializeGetAES(shortKey)
	if err == nil {
		t.Errorf("expected error for incorrect AES key length, got nil")
	}
}

func TestSerializeRedirect(t *testing.T) {
	// Test correct serialization
	validBase64Data := base64.StdEncoding.EncodeToString([]byte("test data"))
	addr := "127.0.0.1:8080"

	rawMsg, err := torpacket.SerialzieRedirect(validBase64Data, addr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check the message code
	if rawMsg.Code != torpacket.Redirect {
		t.Errorf("expected message code %d, got %d", torpacket.Redirect, rawMsg.Code)
	}

	// Check the message content
	var msg torpacket.RedirectMsg
	err = json.Unmarshal([]byte(rawMsg.JsonData), &msg)
	if err != nil {
		t.Errorf("expected no error unmarshalling JSON, got %v", err)
	}

	if msg.Addr != addr {
		t.Errorf("expected address %s, got %s", addr, msg.Addr)
	}

	if msg.RedirectedMessage != validBase64Data {
		t.Errorf("expected data %s, got %s", validBase64Data, msg.RedirectedMessage)
	}

	// Test invalid base64 data
	invalidBase64Data := "invalid_data"
	_, err = torpacket.SerialzieRedirect(invalidBase64Data, addr)
	if err == nil {
		t.Errorf("expected error for invalid base64 data, got nil")
	}
}

func TestSerializeReceive(t *testing.T) {
	// Test correct serialization
	message := "This is a test message"

	rawMsg, err := torpacket.SerializeReceive(message)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check the message code
	if rawMsg.Code != torpacket.Receive {
		t.Errorf("expected message code %d, got %d", torpacket.Receive, rawMsg.Code)
	}

	// Check the message content
	var msg torpacket.ReceiveMsg
	err = json.Unmarshal([]byte(rawMsg.JsonData), &msg)
	if err != nil {
		t.Errorf("expected no error unmarshalling JSON, got %v", err)
	}

	if msg.Message != message {
		t.Errorf("expected message %s, got %s", message, msg.Message)
	}
}

func TestSerializeDestroy(t *testing.T) {
	// Test correct serialization
	rawMsg := torpacket.SerializeDestroy()

	// Check the message code
	if rawMsg.Code != torpacket.Destroy {
		t.Errorf("expected message code %d, got %d", torpacket.Destroy, rawMsg.Code)
	}

	// Check the message content
	if rawMsg.JsonData != "" {
		t.Errorf("expected empty message data, got %s", rawMsg.JsonData)
	}
}

func TestSerializeAck(t *testing.T) {
	// Test correct serialization
	errorMessage := "All good"
	rawMsg, err := torpacket.SerializeAck(errorMessage)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check the message code
	if rawMsg.Code != torpacket.Ack {
		t.Errorf("expected message code %d, got %d", torpacket.Ack, rawMsg.Code)
	}

	// Check the message content
	var msg torpacket.AckMsg
	err = json.Unmarshal([]byte(rawMsg.JsonData), &msg)
	if err != nil {
		t.Errorf("expected no error unmarshalling JSON, got %v", err)
	}

	if msg.Message != errorMessage {
		t.Errorf("expected message %s, got %s", errorMessage, msg.Message)
	}
}

func TestDeserializeMessage(t *testing.T) {
	// Test correct deserialization for GetAES
	aesKey := make([]byte, torpacket.AES_KEY_LENGTH)
	for i := range aesKey {
		aesKey[i] = byte(i)
	}

	rawMsg, _ := torpacket.SerializeGetAES(aesKey)
	parsedMsg, err := torpacket.DeserializeMessage(torpacket.RawMessage{Code: rawMsg.Code, JsonData: rawMsg.JsonData})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check the message code
	if parsedMsg.Code != torpacket.GetAES {
		t.Errorf("expected code %d, got %d", torpacket.GetAES, parsedMsg.Code)
	}

	// Check the deserialized content
	if aesData, ok := parsedMsg.Data.(*torpacket.GetAesMsg); ok {
		expectedData := base64.StdEncoding.EncodeToString(aesKey)
		if aesData.AesKey != expectedData {
			t.Errorf("expected aesKey %s, got %s", expectedData, aesData.AesKey)
		}
	} else {
		t.Errorf("expected *GetAesMsg, got %T", parsedMsg.Data)
	}

	// Test unknown message type
	_, err = torpacket.DeserializeMessage(torpacket.RawMessage{Code: 999, JsonData: "{}"})
	if err == nil {
		t.Errorf("expected error for unknown message type, got nil")
	}

	// Test corrupted JSON data
	_, err = torpacket.DeserializeMessage(torpacket.RawMessage{Code: torpacket.GetAES, JsonData: "{invalid}"})
	if err == nil {
		t.Errorf("expected error for corrupted JSON, got nil")
	}
}
