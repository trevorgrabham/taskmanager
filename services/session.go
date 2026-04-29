package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const UserContextKey = "user_id"
const SessionCookieKey = "session_id"

// 30 days
const SessionTimeoutDuration = 30 * 24 * time.Hour

// GenerateSessionID creates a random 32 byte string.
func GenerateSessionID() (sessionID string) {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)

	sessionID = hex.EncodeToString(bytes)

	return sessionID
}
