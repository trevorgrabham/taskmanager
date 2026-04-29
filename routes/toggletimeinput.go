package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
)

// ToggleTimeHandler renders a TimeInput view.
//
// Requires that the user be authenticated.
func (h Handler) ToggleTimeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                = "ToggleTimeHandler"
		err                   error
		sessionID, timeToggle string
		taskID                int
		t                     services.Task
		user                  services.User
	)
	// Check for an active session
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w %s", caller, ErrInvalidSessionID, sessionID))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if user.ID == 0 {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptySessionID))
		return
	}

	// If it was toggled off, empty response is enough.
	w.Header().Set("Content-Type", "text/html")
	if timeToggle = r.URL.Query().Get("all-day"); timeToggle != "on" {
		fmt.Fprint(w, "")
		return
	}

	if taskID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	// Request from a new task
	if taskID < 1 {
		if err = views.TimeInput("").Render(r.Context(), w); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w TimeInput: %s", caller, ErrRenderingTemplate, err))
			return
		}

		return
	}

	// Request from a task being edited
	if t, err = h.services.GetTaskByID(taskID, user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if err = views.TimeInput(t.DueDate.Format("15:04")).Render(r.Context(), w); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w TimeInput: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
