package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/internal/cookies"
	"local/taskmanager/views"
	"local/taskmanager/views/account"
	"net/http"
)

// SignupHandler renders the Signup view.
// On a POST request, it hands the request off to ParseSignupHandler.
func (h Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.ParseSignupHandler(w, r)
		return
	}

	var (
		caller    = "SignupHandler"
		err       error
		sessionID string
	)

	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	// If no sessionID, render the login form based upon the type of request.
	if sessionID == "" {
		w.Header().Set("Content-Type", "text/html")
		if r.Header.Get("HX-Request") == "true" {
			if err = account.SignupForm().Render(r.Context(), w); err != nil {
				h.HandleError(w, r, fmt.Errorf("%s: %w SignupForm: %s", caller, ErrRenderingTemplate, err))
				return
			}
		} else {
			if err = views.Layout(views.LayoutViewData{UserID: 0, Content: account.SignupForm()}).Render(r.Context(), w); err != nil {
				h.HandleError(w, r, fmt.Errorf("%s: %w Layout: %s", caller, ErrRenderingTemplate, err))
				return
			}
		}

		return
	}

	// If they have a valid session, redirect to home page
	if _, err = h.services.GetUserBySessionID(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w %s", caller, ErrInvalidSessionID, sessionID))
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return
}

// ParseSignupHandler parses the account information from the SignupForm. The data is validated, and upon validation an account is created and a new session is created. The client is given a new SessionID cookie in the response header.
//
// If 'username' is empty, an ErrEmptyUsername is returned.
// If 'password' is empty, an ErrEmptyPassword is returned.
// If 'confirm-password' is empty, an ErrEmptyConfirmPassword is returned.
// If the username exists in the database an ErrUsernameTaken is returned.
// If the password and confirm password do not match, and ErrPasswordsNotMatch is returned.
func (h Handler) ParseSignupHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                                         = "ParseSignupHandler"
		err                                            error
		username, password, confirmPassword, sessionID string
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

	if confirmPassword = r.FormValue("confirm-password"); confirmPassword == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptyConfirmPassword))
		return
	}

	if password != confirmPassword {
		h.HandleError(w, r, fmt.Errorf("%s: %w %s != %s", caller, ErrPasswordsNotMatch, password, confirmPassword))
		return
	}

	if sessionID, err = h.services.SignupUser(username, password); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	cookies.SetSessionIDCookie(w, sessionID)
	w.Header().Set("HX-Redirect", "/")
}
