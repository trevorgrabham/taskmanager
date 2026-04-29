package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/internal/cookies"
	"net/http"
)

// LogoutHandler deletes the session referenced by the clients SessionID cookie, and adds a request header to remove the cookie from the clients browser.
//
// Requires the user to be authenticated.
func (h Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller    = "LogoutHandler"
		err       error
		sessionID string
	)
	// Check user authenticated
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptySessionID))
		return
	}

	if err = h.services.DeleteSession(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	cookies.DeleteSessionIDCookie(w)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
