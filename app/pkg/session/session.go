package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	client *redis.Client
}

type SessionData struct {
	AESKey    string   `json:"aes_key"`
	Addresses []string `json:"addresses"`
}

func NewSessionManager(redisAddr string) *SessionManager {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &SessionManager{
		client: client,
	}
}

// generateSessionToken creates a random session token
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateSession creates a new session with the provided AES key
func (sm *SessionManager) CreateSession(aesKey string) (string, error) {
	ctx := context.Background()

	// Generate session token
	sessionToken, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	// Create session data
	sessionData := SessionData{
		AESKey:    aesKey,
		Addresses: make([]string, 0),
	}

	// Store in Redis with 24 hour expiration
	err = sm.client.HSet(ctx, "session:"+sessionToken, "aes_key", sessionData.AESKey).Err()
	if err != nil {
		return "", err
	}

	err = sm.client.Expire(ctx, "session:"+sessionToken, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

// AppendAddr adds an address to the session
func (sm *SessionManager) AppendAddr(sessionToken, addr string) error {
	ctx := context.Background()

	// Check if session exists
	exists, err := sm.client.Exists(ctx, "session:"+sessionToken).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("session not found")
	}

	// Append address to the list
	err = sm.client.SAdd(ctx, "session:"+sessionToken+":addresses", addr).Err()
	return err
}

// PullData retrieves all session data
func (sm *SessionManager) PullData(sessionToken string) (*SessionData, error) {
	ctx := context.Background()

	// Check if session exists
	exists, err := sm.client.Exists(ctx, "session:"+sessionToken).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, errors.New("session not found")
	}

	// Get AES key
	aesKey, err := sm.client.HGet(ctx, "session:"+sessionToken, "aes_key").Result()
	if err != nil {
		return nil, err
	}

	// Get addresses
	addresses, err := sm.client.SMembers(ctx, "session:"+sessionToken+":addresses").Result()
	if err != nil {
		return nil, err
	}

	return &SessionData{
		AESKey:    aesKey,
		Addresses: addresses,
	}, nil
}
