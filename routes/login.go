package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/internal/cookies"
	"local/taskmanager/views/account"
	"net/http"
)

func (h Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.ParseLoginHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := account.LoginForm().Render(r.Context(), w); err != nil {
		h.HandleError(w, NewTemplateRenderError("LoginHandler", "LoginForm", err))
		return
	}
}

func (h Handlers) ParseLoginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                        = "ParseLoginHandler"
		err                           error
		sessionID, username, password string
	)

	if sessionID, err = context.GetSessionID(r.Context()); sessionID != "" && err == nil {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	if err = r.ParseForm(); err != nil {
		h.HandleError(w, NewParseFormError(caller, err))
		return
	}

	if username = r.FormValue("username"); username == "" {
		h.HandleError(w, NewNoUsernameError(caller))
		return
	}

	if password = r.FormValue("password"); password == "" {
		h.HandleError(w, NewNoPasswordError(caller))
		return
	}

	if sessionID, err = h.services.LoginUser(username, password); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	cookies.SetSessionIDCookie(w, sessionID)
	w.Header().Set("HX-Redirect", "/")
}
