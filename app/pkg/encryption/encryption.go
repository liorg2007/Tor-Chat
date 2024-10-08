package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
)

const (
	AES_KEY_SIZE = 128
	RSA_KEY_SIZE = 2048
)

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

func (r *RSAEncryptor) GenerateKey() error {
	privKey, err := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)

	if err != nil {
		return err
	}

	r.PrivateKey = privKey
	r.PublicKey = &privKey.PublicKey

	return nil
}

func (a *AESEncryptor) GenerateKey() error {
	a.Key = make([]byte, AES_KEY_SIZE)
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
		return nil, err
	}

	return ciphertext, nil
}

func (a *AESEncryptor) Encrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, text, nil)

	return ciphertext, nil
}

func (r *RSAEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	hash := sha256.New()
	text, err := rsa.DecryptOAEP(hash, rand.Reader, r.PrivateKey, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return text, nil
}

func (a *AESEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedData, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
