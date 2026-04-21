package account

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

const UserContextKey = "user_id"
const SessionCookieKey = "session_id"

// 30 days
const SessionTimeoutDuration = 30 * 24 * time.Hour 

func GenerateSessionID() (sessionID string, err error) {
	bytes := make([]byte, 32)
	if _, err = rand.Read(bytes); err != nil { return "", fmt.Errorf("CreateSession: %s", err) }

	sessionID = hex.EncodeToString(bytes)

	return sessionID, nil
}
