package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
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
}

type AESEncryptor struct {
	Key []byte
}

func (r *RSAEncryptor) GenerateKey(size int) error {
	privKey, err := rsa.GenerateKey(rand.Reader, size)

	if err != nil {
		log.Panic(err)
		return err
	}

	r.PrivateKey = privKey
	r.PublicKey = &privKey.PublicKey

	return nil
}

func (a *AESEncryptor) GenerateKey(size int) error {
	a.Key = make([]byte, size)
	_, err := rand.Read(a.Key)
	if err != nil {
		return err
	}
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
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Panic("error generating the nonce ", err)
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, text, nil)

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
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	decryptedData, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return decryptedData, nil
}
