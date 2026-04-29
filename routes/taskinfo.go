package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views/taskinfo"
	"net/http"
)

// TaskInfoHandler renders the TaskInfo view for the task identified by 'id'.
//
// Requires that the user be authenticated.
func (h Handler) TaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller       = "TaskInfoHandler"
		err          error
		sessionID    string
		user         services.User
		taskID       int
		t            services.Task
		taskInfoData taskinfo.TaskInfoViewData
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

	if t, err = h.services.GetTaskByID(taskID, user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	taskInfoData = h.parseTaskToTaskInfoViewData(t)

	w.Header().Set("Content-Type", "text/html")
	err = taskinfo.TaskInfo(taskInfoData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w TaskInfo: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
