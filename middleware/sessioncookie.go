// Package middleware provides cookie-parsing middleware for populating http.Request contexts.
package middleware

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/internal/cookies"
	"net/http"
)

// SessionIDCookieKey is the default key for SessionID cookies.
const SessionIDCookieKey = "session_id"

// SessionIDCookieMiddleware populates the requests context if a SessionID cookie is present and valid. If the cookie is absent the request proceeds without a SessionID in the context. If the cookie is present, but not valid, a response header is written to request the removal of the cookie.
func SessionIDCookieMiddleware(next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			sessionCookie *http.Cookie
			err           error
		)
		if sessionCookie, err = r.Cookie(SessionIDCookieKey); err != nil {
			// No cookie found, pass through
			next.ServeHTTP(w, r)
			return
		}

		if err = sessionCookie.Valid(); err != nil {
			// Malformed cookie, clear it and pass through
			cookies.DeleteSessionIDCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.SetSessionID(r.Context(), sessionCookie.Value)))
	}
}
