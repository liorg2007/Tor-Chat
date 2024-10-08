package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// Encode RSA public key to base64
func EncodeRSAPublicKey(publicKey *rsa.PublicKey) (string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}
	return base64.StdEncoding.EncodeToString(derBytes), nil
}

// Decode RSA public key from base64
func DecodeRSAPublicKey(encodedKey string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}
	parsedKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPublicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPublicKey, nil
}

// Encode AES key (byte array) to base64
func EncodeAESKey(aesKey []byte) string {
	return base64.StdEncoding.EncodeToString(aesKey)
}

// Decode AES key from base64 to byte array
func DecodeAESKey(encodedKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedKey)
}
