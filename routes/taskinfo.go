package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views/taskinfo"
	"net/http"
)

// Renders task-info page
func (h Handlers) TaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller       = "TaskInfoHandler"
		err          error
		sessionID    string
		user         services.User
		taskID       int
		t            services.Task
		taskInfoData taskinfo.TaskInfoViewData
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

	if t, err = h.services.GetTaskByID(taskID, user.ID); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	taskInfoData = h.parseTaskToTaskInfoViewData(t)

	w.Header().Set("Content-Type", "text/html")
	err = taskinfo.TaskInfo(taskInfoData).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "TaskInfo", err))
		return
	}
}
