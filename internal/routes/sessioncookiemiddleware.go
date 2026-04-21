package routes

import (
	"context"
	"local/taskmanager/internal/account"
	sqlite "local/taskmanager/internal/db"
	"net/http"
)

func (h Handlers) SessionCookieMiddleware(next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			sessionCookie *http.Cookie
			err           error
			userID        int
			ctx           context.Context
		)
		if sessionCookie, err = r.Cookie(account.SessionCookieKey); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if userID, err = sqlite.GetUserID(h.DB, sessionCookie.Value); err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session_id", MaxAge: -1})
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(r.Context(), account.UserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
