package routes

import (
	"local/taskmanager/internal/views/user"
	"log"
	"net/http"
)

func (h Handlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.ParseSignupHandler(w, r)
		return
	}
	var ok bool
	if ok = h.checkConnection(w); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := user.SignupForm().Render(r.Context(), w); err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("SignupHandler: %s\n", err)
		return
	}
}

func (h Handlers) ParseSignupHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok                                  bool
		username, password, confirmPassword string
	)
	if ok = h.checkConnection(w); !ok {
		return
	}

	if ok = h.parseForm(w, r); !ok {
		return
	}

	username = r.FormValue("username")
	password = r.FormValue("password")
	confirmPassword = r.FormValue("confirm-password")

	if password != confirmPassword {
		http.Error(w, "Error passwords do not match", http.StatusBadRequest)
		log.Printf("SignUpHandler: password != confirm-password (%s != %s)\n", password, confirmPassword)
		return
	}

	if ok = h.addUser(w, username, password); !ok {
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
