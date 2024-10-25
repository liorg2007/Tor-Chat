package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
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

// Exposed function: Encrypts using RSA and returns base64-encoded ciphertext
func (r *RSAEncryptor) EncryptBase64(textBase64 string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(textBase64)
	if err != nil {
		return "", err
	}

	ciphertext, err := r.Encrypt(text)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Exposed function: Decrypts base64-encoded ciphertext using RSA and returns base64-encoded plaintext
func (r *RSAEncryptor) DecryptBase64(ciphertextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	plaintext, err := r.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(plaintext), nil
}

// Exposed function: Encrypts using AES and returns base64-encoded ciphertext
func (a *AESEncryptor) EncryptBase64(textBase64 string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(textBase64)
	if err != nil {
		return "", err
	}

	ciphertext, err := a.Encrypt(text)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Exposed function: Decrypts base64-encoded ciphertext using AES and returns base64-encoded plaintext
func (a *AESEncryptor) DecryptBase64(ciphertextBase64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	plaintext, err := a.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(plaintext), nil
}
