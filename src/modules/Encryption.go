package Encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
)

type Encrypor interface {
	Encrypt(text []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	GenerateKey(size int) error
}

type RSAEncryptor struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	size       int
}

type AESEncryptor struct {
	Key  []byte
	size int
}

func (r *RSAEncryptor) GenerateKey(size int) error {
	privKey, err := rsa.GenerateKey(rand.Reader, size)

	if err != nil {
		log.Panic(err)
		return err
	}

	r.PrivateKey = privKey
	r.PublicKey = &privKey.PublicKey
	r.size = size

	return nil
}

func (a *AESEncryptor) GenerateKey(size int) error {
	a.Key = make([]byte, size)
	a.size = size

	return nil
}
