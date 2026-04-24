package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
)

func (h Handlers) ToggleTimeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                            = "ToggleTimeHandler"
		err                               error
		sessionID, timeToggle, timeString string
		t                                 services.Task
		user                              services.User
	)
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, NewSessionIDCookieParseError(caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, NewUnauthenticatedError(caller))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if timeToggle = r.URL.Query().Get("all-day"); timeToggle != "on" {
		fmt.Fprint(w, "")
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

	if t.ID > 0 {
		if t, err = h.services.GetTaskByID(t.ID, user.ID); err != nil {
			// TODO: wrap the returned error
			h.HandleError(w, err)
			return
		}
	}

	timeString = t.DueDate.Format("15:04")

	err = views.TimeInput(timeString).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "TimeInput", err))
		return
	}
}
