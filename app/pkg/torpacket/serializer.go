package torpacket

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"marshmello/pkg/helper"
)

const AES_KEY_LENGTH int = 32 // 256 bits = 32 bytes

// SerializeGetAES serializes a message containing an AES key.
// The AES key is base64 encoded before being included in the message.
// It returns a RawMessage with the serialized data or an error if the key length is invalid.
//
// Parameters:
//   - aesKey: A byte slice representing the AES key (must be 32 bytes).
//
// Returns:
//   - RawMessage: The serialized message containing the AES key.
//   - error: Error if the AES key length is not exactly 32 bytes.
func SerializeGetAES(aesKey []byte) (RawMessage, error) {
	if len(aesKey) != AES_KEY_LENGTH {
		return RawMessage{}, errors.New("aes Key is not in desired length")
	}

	m := GetAesMsg{AesKey: base64.StdEncoding.EncodeToString(aesKey)}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{GetAES, string(b)}, nil
}

// SerialzieRedirect serializes a message to redirect traffic to another node.
// The encrypted message should be base64 encoded and will be sent to the provided address.
//
// Parameters:
//   - encryptedData: A base64-encoded string representing the encrypted data to be redirected.
//   - addr: A string representing the address to which the message should be redirected.
//
// Returns:
//   - RawMessage: The serialized redirect message.
//   - error: Error if the encrypted data is not base64 encoded.
func SerialzieRedirect(encryptedData string, addr string) (RawMessage, error) {
	if !helper.IsBase64Encoded(encryptedData) {
		return RawMessage{}, errors.New("encrypted data is not base64 encoded")
	}

	m := RedirectMsg{Addr: addr, RedirectedMessage: encryptedData}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Redirect, string(b)}, nil
}

// SerializeReceive serializes a message to be sent to a node.
// This message contains arbitrary data that will be processed by the receiver.
//
// Parameters:
//   - message: A string representing the message to be sent.
//
// Returns:
//   - RawMessage: The serialized message for receiving data.
//   - error: Error if the message could not be serialized.
func SerializeReceive(message string) (RawMessage, error) {
	m := ReceiveMsg{message}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Receive, string(b)}, nil
}

// SerializeDestroy serializes a message indicating that the circuit or session should be destroyed.
// This message typically signals the end of communication between nodes.
//
// Returns:
//   - RawMessage: The serialized destroy message (no payload).
func SerializeDestroy() RawMessage {
	return RawMessage{Destroy, ""}
}

// SerializeAck serializes an acknowledgment message to confirm the receipt of data.
// This message is often used to notify the sender that a message has been successfully received and processed.
//
// Parameters:
//   - errorMessage: A string representing the acknowledgment message or error to be returned.
//
// Returns:
//   - RawMessage: The serialized acknowledgment message.
//   - error: Error if the message could not be serialized.
func SerializeAck(errorMessage string) (RawMessage, error) {
	m := AckMsg{errorMessage}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Ack, string(b)}, nil
}
