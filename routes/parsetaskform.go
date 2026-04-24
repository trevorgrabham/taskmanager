package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views/taskinfo"
	"log"
	"net/http"
)

func (h Handlers) ParseTaskForm(w http.ResponseWriter, r *http.Request) {
	var (
		caller                           = "ParseTaskForm"
		err                              error
		sessionID, dayString, timeString string
		user                             services.User
		t                                services.Task
		taskInfoData                     taskinfo.TaskInfoViewData
	)
	if r.Method != http.MethodPost {
		h.HandleError(w, NewWrongMethodError(caller, http.MethodPost, r.Method))
		return
	}

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

	t.UserID = user.ID

	if err = r.ParseForm(); err != nil {
		h.HandleError(w, NewParseFormError(caller, err))
		return
	}

	if t.Title = r.FormValue("title"); t.Title == "" {
		h.HandleError(w, NewNoTitleError(caller))
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

	t.Category = r.FormValue("category") // TODO: handle uncatagorized tasks here. Potentially error
	t.Description = r.FormValue("description")
	dayString = r.FormValue("due-date")
	timeString = r.FormValue("time")
	if t.DueDate, err = h.parseDateAndTime(dayString, timeString); err != nil {
		h.HandleError(w, NewDateTimeParseError(caller, fmt.Sprintf("%s %s", dayString, timeString), err))
		return
	}

	if t.RecurringPeriod, err = h.parseRecurringPeriod(r.FormValue("due-date"), r.FormValue("time")); err != nil {
		h.HandleError(w, NewRecurringParseError(caller, err))
		return
	}

	switch {
	case t.ID == 0:
		if t, err = h.services.AddTask(t, user.ID); err != nil {
			// TODO: wrap the returned error
			h.HandleError(w, err)
			return
		}

		w.Header().Set("HX-Redirect", "/")
	case t.ID > 1:
		if t, err = h.services.UpdateTask(t, user.ID); err != nil {
			// TODO: wrap the returned error
			h.HandleError(w, err)
			return
		}

		taskInfoData = h.parseTaskToTaskInfoViewData(t)

		w.Header().Set("Content-Type", "text/html")
		err := taskinfo.TaskInfo(taskInfoData).Render(r.Context(), w)
		if err != nil {
			h.HandleError(w, NewTemplateRenderError(caller, "TaskInfo", err))
			return
		}
	default:
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("parsing task form: bad id value of %d\n", t.ID)
		return
	}
}
