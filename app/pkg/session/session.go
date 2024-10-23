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
	AESKey  string `json:"aes_key"`
	Address string `json:"address"`
}

func NewSessionManager(redisAddr string) (*SessionManager, error) {
	err := error(nil)
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Create context for Redis operations
	ctx := context.Background()

	// Test the connection
	_, err = client.Ping(ctx).Result()

	return &SessionManager{
		client: client,
	}, err
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
		AESKey:  aesKey,
		Address: "",
	}

	// Store in Redis with 24-hour expiration
	err = sm.client.HSet(ctx, "session:"+sessionToken, "aes_key", sessionData.AESKey).Err()
	if err != nil {
		return "", err
	}

	err = sm.client.HSet(ctx, "session:"+sessionToken, "address", sessionData.Address).Err()
	if err != nil {
		return "", err
	}

	err = sm.client.Expire(ctx, "session:"+sessionToken, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

// UpdateAddress updates the address in the session
func (sm *SessionManager) UpdateAddress(sessionToken, addr string) error {
	ctx := context.Background()

	// Check if session exists
	exists, err := sm.client.Exists(ctx, "session:"+sessionToken).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("session not found")
	}

	// Update address
	err = sm.client.HSet(ctx, "session:"+sessionToken, "address", addr).Err()
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

	// Get address
	address, err := sm.client.HGet(ctx, "session:"+sessionToken, "address").Result()
	if err != nil {
		return nil, err
	}

	return &SessionData{
		AESKey:  aesKey,
		Address: address,
	}, nil
}
