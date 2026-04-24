package context

import (
	"context"
	"fmt"
)

type contextKey struct{}
var sessionIDContextKey = contextKey{}

func GetSessionID(ctx context.Context) (sessionID string, err error) {
	contextValue := ctx.Value(sessionIDContextKey)
	if contextValue == nil {
		return "", nil
	}

	var ok bool
	if sessionID, ok = contextValue.(string); !ok {
		return "", fmt.Errorf("cannot convert %v into a string value", contextValue)
	}

	return sessionID, nil
}

func SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDContextKey, sessionID)
}
