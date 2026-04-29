package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
	"time"
)

// AddTaskHandler renders the TaskForm view.
//
// If a 'due-date' query string parameter is present, it is parsed and used to fill out the 'due-date' input for the form. If 'due-date' is not a valid date-time string, then it is ignored and an ErrInvalidDueDate is logged.
// If a 'category' query string parameter is present, it is used to fill out the 'category' input for the form.
//
// Requires that the user be authenticated.
func (h Handler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller              = "AddTaskHandler"
		err                 error
		sessionID           string
		taskViewData        views.TaskFormViewData
		user                services.User
		dueDay              time.Time
		categorySuggestions []string
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

	if categorySuggestions, err = h.services.GetUserCategorySuggestions(user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if dueDay, err = h.parseDay(r.URL.Query().Get("due-date")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if dueDay.IsZero() {
		taskViewData.DueDateDay = ""
	} else {
		taskViewData.DueDateDay = dueDay.Format("2006-01-02")
	}
	taskViewData.Category = r.URL.Query().Get("category")
	taskViewData.UserID = user.ID
	taskViewData.CategorySuggestions = categorySuggestions

	w.Header().Set("Content-Type", "text/html")
	err = views.TaskForm(taskViewData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w TaskForm: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
