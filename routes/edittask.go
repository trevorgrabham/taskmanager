package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
)

// GET request: Renders task-info page in view mode
func (h Handlers) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if taskID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, NewIDParseError(caller, err))
		return
	}
	if taskID == -1 {
		h.HandleError(w, NewNoIDError(caller))
		return
	}

	if taskData, err = h.services.GetTaskByID(taskID, user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if taskViewData, err = h.parseTaskToTaskFormViewData(taskData); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if categorySuggestions, err = h.services.GetUserCategorySuggestions(user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	taskViewData.CategorySuggestions = categorySuggestions

	w.Header().Set("Content-Type", "text/html")
	err = views.TaskForm(taskViewData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "TaskForm", err))
		return
	}
}
