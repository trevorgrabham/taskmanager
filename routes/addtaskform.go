package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
	"time"
)

// GET: Renders the add-task form
func (h Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller              = "AddTaskHandler"
		err                 error
		sessionID           string
		taskViewData        views.TaskFormViewData
		user                services.User
		dueDay              time.Time
		categorySuggestions []string
	)

	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, NewSessionIDCookieParseError(caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, NewUnauthenticatedError(caller))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if categorySuggestions, err = h.services.GetUserCategorySuggestions(user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if dueDay, err = h.parseDay(r.URL.Query().Get("due-date")); err != nil {
		h.HandleError(w, NewDayParseError(caller, err))
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
		h.HandleError(w, NewTemplateRenderError(caller, "TaskForm", err))
		return
	}
}
