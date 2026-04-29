// Package context provides functions for storing and retrieving values using context.Context objects.
package context

import (
	"context"
	"errors"
	"fmt"
)

type contextKey struct{}

var sessionIDContextKey = contextKey{}
var (
	ErrWrongContextType = errors.New("wrong context type")
)

// GetSessionID retrieves the SessionID string from the provided context.Context.
//
// If no SessionID is stored with the context, then "" is returned.
// Returns an ErrWrongContextType if a SessionID is stored, but it is not a valid string.
func GetSessionID(ctx context.Context) (sessionID string, err error) {
	contextValue := ctx.Value(sessionIDContextKey)
	if contextValue == nil {
		return "", nil
	}

	var ok bool
	if sessionID, ok = contextValue.(string); !ok {
		return "", fmt.Errorf("%w for sessionID: %s", ErrWrongContextType, contextValue)
	}

	return sessionID, nil
}

// SetSessionID stores a SessionID string on the given context.Context and returns it.
func SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDContextKey, sessionID)
}
