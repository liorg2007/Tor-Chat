package Encryption

import "crypto/rsa"

type Encrypor interface {
	Encrypt(text []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	GenerateKey() error
}

type RSAEncryptor struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type AESEncryptor struct {
	Key []byte
}
