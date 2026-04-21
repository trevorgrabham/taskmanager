package routes

import (
	"local/taskmanager/internal/account"
	"local/taskmanager/internal/views/user"
	"log"
	"net/http"
	"time"
)

func (h Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.ParseLoginHandler(w, r)
		return
	}
	var ok bool
	if ok = h.checkConnection(w); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := user.LoginForm().Render(r.Context(), w); err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("LoginHandler: %s\n", err)
		return
	}
}

func (h Handlers) ParseLoginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok                            bool
		username, password, sessionID string
		user                          account.User
	)
	if ok = h.checkConnection(w); !ok {
		return
	}

	if ok = h.parseForm(w, r); !ok {
		return
	}

	username = r.FormValue("username")
	password = r.FormValue("password")

	if user, ok = h.getUserByUsername(w, username); !ok {
		return
	}

	if ok = h.comparePassword(w, user.HashedPassword, password); !ok {
		return
	}

	if sessionID, ok = h.startSession(w, user.ID); !ok {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     account.SessionCookieKey,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(account.SessionTimeoutDuration),
	})

	w.Header().Set("HX-Redirect", "/")
}
