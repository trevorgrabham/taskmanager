package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
)

// ToggleCompleteHandler renders a ListItem view with the updated completion status of the task.
//
// Requires that the user be authenticated.
func (h Handler) ToggleCompleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller       = "ToggleCompleteHandler"
		err          error
		sessionID    string
		user         services.User
		taskID int
		t            services.Task
		taskViewData views.TaskViewData
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

	if taskID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if t, err = h.services.TaskToggleComplete(taskID, user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	taskViewData = h.parseTaskToTaskViewData(t)

	w.Header().Set("Content-Type", "text/html")
	err = views.ListItem(taskViewData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w ListItem: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
