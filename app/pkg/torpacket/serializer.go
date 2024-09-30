package torpacket

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"marshmello/pkg/helper"
)

const AES_KEY_LENGTH int = 256

const (
	GetAES   = 100 // Create a new circuit
	Redirect = 101
	Receive  = 102 // Custom: Redirect to another server
	Destroy  = 103 // Custom: Receive data
	Ack      = 104
)

func SerializeGetAES(aesKey []byte) (RawMessage, error) {
	if len(aesKey) != AES_KEY_LENGTH {
		return RawMessage{}, errors.New("Aes Key is not in desired length")
	}

	m := GetAesMsg{AesKey: base64.StdEncoding.EncodeToString(aesKey)}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{GetAES, string(b)}, nil
}

func SerialzieRedirect(encryptedData string, addr string) (RawMessage, error) {
	if helper.IsBase64Encoded(encryptedData) {
		return RawMessage{}, errors.New("encrypted data is not base64")
	}

	m := RedirectMsg{Addr: addr, RedirectedMessage: encryptedData}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Redirect, string(b)}, nil
}

func SerializeReceive(message string) (RawMessage, error) {
	m := ReceiveMsg{message}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Receive, string(b)}, nil
}

func SerializeDestroy() (RawMessage, error) {
	return RawMessage{Destroy, ""}, nil
}

func SerializeAck(errorMessage string) (RawMessage, error) {
	m := AckMsg{errorMessage}

	b, err := json.Marshal(m)

	if err != nil {
		return RawMessage{}, err
	}

	return RawMessage{Destroy, string(b)}, nil
}
