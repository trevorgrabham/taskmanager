package cookies

import (
	"net/http"
	"time"
)

const SessionCookieKey = "session_id"
const SessionTimeoutDuration = 30 * 24 * time.Hour

func SetSessionIDCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieKey,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,			// TODO: switch to 'true' once https enabled
		SameSite: http.SameSiteLaxMode,		// Look into swapping this once I don't need to debug using chrome dev tools anymore
		Expires:  time.Now().Add(SessionTimeoutDuration),
	})
}

func GetSessionIDCookie() {
	// update middleware after implemented
}
