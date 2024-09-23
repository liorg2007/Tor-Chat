package Encryption

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

func (r *RSAEncryptor) Encrypt(text []byte) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, r.PublicKey, text, nil)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return ciphertext, nil
}

func (a *AESEncryptor) Encrypt(text []byte) ([]byte, error) {
	aes, err := aes.NewCipher(a.Key)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	ciphertext := make([]byte, len(text))
	aes.Encrypt(ciphertext, text)

	return ciphertext, nil
}

func (r *RSAEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	hash := sha256.New()
	text, err := rsa.DecryptOAEP(hash, rand.Reader, r.PrivateKey, ciphertext, nil)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return text, nil
}

func (a *AESEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	aes, err := aes.NewCipher(a.Key)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	aes.Decrypt(plaintext, ciphertext)

	return plaintext, nil
}
