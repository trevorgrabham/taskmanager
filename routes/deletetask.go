package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"net/http"
)

// Hits the backend db.DeleteTask() and redirects to "/"
func (h Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller    = "DeleteTaskHandler"
		err       error
		sessionID string
		taskID    int
		user      services.User
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

	if err = h.services.DeleteTask(taskID, user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
