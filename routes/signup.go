package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/internal/cookies"
	user "local/taskmanager/views/account"
	"net/http"
)

func (h Handlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.ParseSignupHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := user.SignupForm().Render(r.Context(), w); err != nil {
		h.HandleError(w, NewTemplateRenderError("SignupHandler", "SignupForm", err))
		return
	}
}

func (h Handlers) ParseSignupHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                                         = "ParseSignupHandler"
		err                                            error
		username, password, confirmPassword, sessionID string
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

	if confirmPassword = r.FormValue("confirm-password"); confirmPassword == "" {
		h.HandleError(w, NewNoConfirmPasswordError(caller))
		return
	}

	if password != confirmPassword {
		h.HandleError(w, NewMismatchPasswordError(caller, password, confirmPassword))
		return
	}

	if sessionID, err = h.services.SignupUser(username, password); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	cookies.SetSessionIDCookie(w, sessionID)
	w.Header().Set("HX-Redirect", "/")
}
