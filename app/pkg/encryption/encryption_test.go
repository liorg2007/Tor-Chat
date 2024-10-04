package encryption

import (
	"bytes"
	"testing"
)

func TestRSAEncryption(t *testing.T) {
	// Create an RSA encryptor
	rsaEncryptor := &RSAEncryptor{}
	err := rsaEncryptor.GenerateKey(2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Test data to encrypt
	plaintext := []byte("This is a test message.")

	// Encrypt the plaintext
	ciphertext, err := rsaEncryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Decrypt the ciphertext
	decryptedText, err := rsaEncryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare decrypted text with the original plaintext
	if !bytes.Equal(decryptedText, plaintext) {
		t.Errorf("Decrypted text does not match the original. Expected %s, got %s", plaintext, decryptedText)
	}
}

func TestAESEncryption(t *testing.T) {
	// Create an AES encryptor
	aesEncryptor := &AESEncryptor{}
	err := aesEncryptor.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}

	// Test data to encrypt
	plaintext := []byte("This is a test message.")

	// Encrypt the plaintext
	ciphertext, err := aesEncryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Decrypt the ciphertext
	decryptedText, err := aesEncryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare decrypted text with the original plaintext
	if !bytes.Equal(decryptedText, plaintext) {
		t.Errorf("Decrypted text does not match the original. Expected %s, got %s", plaintext, decryptedText)
	}
}
