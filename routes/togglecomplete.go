package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"local/taskmanager/views/dashboard"
	"net/http"
)

func (h Handlers) ToggleCompleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller    = "ToggleCompleteHandler"
		err       error
		sessionID string
		user services.User
		t         services.Task
		taskViewData views.TaskViewData
	)
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, NewSessionIDCookieParseError(caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, NewUnauthenticatedError(caller))
		return
	}

	if t.ID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, NewIDParseError(caller, err))
		return
	}
	if t.ID == -1 {
		h.HandleError(w, NewNoIDError(caller))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if t, err = h.services.TaskToggleComplete(t.ID, user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	taskViewData = h.parseTaskToTaskViewData(t) 

	w.Header().Set("Content-Type", "text/html")
	err = dashboard.ListItem(taskViewData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "ListItem", err))
		return
	}
}
