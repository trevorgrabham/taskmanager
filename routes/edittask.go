package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
)

// EditTaskHandler renders the EditPage view for the task matching 'id'.
//
// If id is empty, an ErrInvalidID is returned.
//
// Requires that the user be authenticated.
func (h Handler) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller              = "EditTaskHandler"
		err                 error
		sessionID           string
		taskID              int
		user                services.User
		taskData            services.Task
		taskViewData        views.TaskFormViewData
		categorySuggestions []string
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

	if taskData, err = h.services.GetTaskByID(taskID, user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if taskViewData, err = h.parseTaskToTaskFormViewData(taskData); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if categorySuggestions, err = h.services.GetUserCategorySuggestions(user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	taskViewData.CategorySuggestions = categorySuggestions

	w.Header().Set("Content-Type", "text/html")
	err = views.TaskForm(taskViewData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w TaskForm: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
