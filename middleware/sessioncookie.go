package middleware

import (
	"local/taskmanager/internal/context"
	"net/http"
)

const SessionIDCookieKey = "session_id"

// Adds the SessionID to r.Context() if a SessionID cookie:
// a) exists
// b) is valid
func SessionIDCookieMiddleware(next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			sessionCookie *http.Cookie
			err           error
		)
		if sessionCookie, err = r.Cookie(SessionIDCookieKey); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if err = sessionCookie.Valid(); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.SetSessionID(r.Context(), sessionCookie.Value)))
	}
}
