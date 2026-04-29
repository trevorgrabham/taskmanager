package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/internal/cookies"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"local/taskmanager/views/account"
	"net/http"
)

// LoginHandler renders the Login view.
// On a POST request, it hands the request off to ParseLoginHandler.
func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// If it's a post request, then we should parse the login form
	if r.Method == http.MethodPost {
		h.ParseLoginHandler(w, r)
		return
	}

	var (
		caller    = "LoginHandler"
		sessionID string
		err       error
		user      services.User
	)

	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	// If no sessionID, render the login form based upon the type of request.
	if sessionID == "" {
		w.Header().Set("Content-Type", "text/html")
		if r.Header.Get("HX-Request") == "true" {
			if err = account.LoginForm().Render(r.Context(), w); err != nil {
				h.HandleError(w, r, fmt.Errorf("%s: %w LoginForm: %s", caller, ErrRenderingTemplate, err))
				return
			}
		} else {
			if err = views.Layout(views.LayoutViewData{ UserID: 0, Content: account.LoginForm()}).Render(r.Context(), w); err != nil {
				h.HandleError(w, r, fmt.Errorf("%s: %w Layout: %s", caller, ErrRenderingTemplate, err))
				return
			}
		}

		return
	}

	// If they have a valid session, redirect to home page
	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w %s", caller, ErrInvalidSessionID, sessionID))
		return
	}
	if user.ID == 0 {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptySessionID))
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return
}


// ParseLoginHandler parses the account information from the LoginForm. The data is validated, and upon validation a new session is created. The client is given a new SessionID cookie in the response header.
//
// If 'username' is empty, an ErrEmptyUsername is returned.
// If 'password' is empty, an ErrEmptyPassword is returned.
// If the username does not exist in the database an ErrUsernameNotExist is returned.
// If the password does not match an ErrWrongPassword is returned.
func (h Handler) ParseLoginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                        = "ParseLoginHandler"
		err                           error
		sessionID, username, password string
	)

	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	// If they have a valid session, redirect them to the home page.
	if sessionID != "" {
		if _, err = h.services.GetUserBySessionID(sessionID); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
			return
		}
		w.Header().Set("HX-Redirect", "/")
		return
	}

	if err = r.ParseForm(); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrParsingForm))
		return
	}

	if username = r.FormValue("username"); username == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptyUsername))
		return
	}

	if password = r.FormValue("password"); password == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptyPassword))
		return
	}

	if sessionID, err = h.services.LoginUser(username, password); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	cookies.SetSessionIDCookie(w, sessionID)
	w.Header().Set("HX-Redirect", "/")
}
