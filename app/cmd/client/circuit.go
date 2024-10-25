package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"marshmello/pkg/encryption"
	"marshmello/pkg/handlers"
)

// GetAesKey requests the AES key and session token from the server
func GetAesKey(addr string) ([]byte, string, error) {
	var req handlers.GetAesRequest
	var rsaEnc encryption.RSAEncryptor
	var res handlers.GetAesResponse

	// Generate RSA keys
	err := rsaEnc.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	// Encode the RSA public key for transmission
	req.RsaKey, err = encryption.EncodeRSAPublicKey(rsaEnc.PublicKey)
	if err != nil {
		return nil, "", err
	}

	// Send request to /get-aes
	respData, err := SendHttpRequest(addr, req, "get-aes")
	if err != nil {
		return nil, "", err
	}

	// Unmarshal the response
	if err := json.Unmarshal(respData, &res); err != nil {
		return nil, "", errors.New("error decoding AES response")
	}

	// Decode the base64-encoded AES key from the response
	aesKey, err := base64.StdEncoding.DecodeString(res.Aes_key)
	if err != nil {
		return nil, "", errors.New("error decoding AES key")
	}

	return aesKey, res.Session, nil
}
