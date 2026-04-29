package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"net/http"
)

// DeleteTaskHandler deletes the task identified by 'id' if it is owned by the user.
//
// If id is empty, an ErrInvalidID is returned.
//
// Requires that the user be authenticated.
func (h Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller    = "DeleteTaskHandler"
		err       error
		sessionID string
		taskID    int
		user      services.User
	)
	// Check user authenticated and grab User
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptySessionID))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	// Get taskID
	if taskID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if err = h.services.DeleteTask(taskID, user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
