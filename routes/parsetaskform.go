package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views/taskinfo"
	"log"
	"net/http"
)

// ParseTaskForm parses values from the TaskForm for adding or updating tasks.
//
// If 'id' is empty, then the parsed data is used to add a new task. If 'id' is present then the data is used to updatethe task using 'id' as an identifier.
// If 'title' is empty, a warning is returned to the user and an error is logged.
// If 'due-date' is not a valid date-time string, an error is logged.
//
// Requires that the user be authenticated.
func (h Handler) ParseTaskForm(w http.ResponseWriter, r *http.Request) {
	var (
		caller       = "ParseTaskForm"
		err          error
		sessionID    string
		user         services.User
		t            services.Task
		taskInfoData taskinfo.TaskInfoViewData
	)
	if r.Method != http.MethodPost {
		h.HandleError(w, r, fmt.Errorf("%s: %w got: %s, wanted: %s", caller, ErrWrongMethod, r.Method, http.MethodPost))
		return
	}

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

	// Parse the form
	if err = r.ParseForm(); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrParsingForm))
		return
	}

	if t.ID, err = h.parseID(r.FormValue("id")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if t.Title = r.FormValue("title"); t.Title == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptyTitle))
		return
	}

	t.UserID = user.ID
	t.Category = r.FormValue("category") // TODO: handle uncatagorized tasks here. Potentially error
	t.Description = r.FormValue("description")
	if t.DueDate, err = h.parseDateAndTime(r.FormValue("due-date"), r.FormValue("time")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if t.RecurringPeriod, err = h.parseRecurringPeriod(r.FormValue("period-value"), r.FormValue("period-unit")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	switch {
	case t.ID == 0:
		if t, err = h.services.AddTask(t); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
			return
		}

		w.Header().Set("HX-Redirect", "/")
	case t.ID > 1:
		if t, err = h.services.UpdateTask(t); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
			return
		}

		taskInfoData = h.parseTaskToTaskInfoViewData(t)

		w.Header().Set("Content-Type", "text/html")
		err := taskinfo.TaskInfo(taskInfoData).Render(r.Context(), w)
		if err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w TaskInfo: %s", caller, ErrRenderingTemplate, err))
			return
		}
	default:
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("%s: bad id value of %d\n", caller, t.ID)
		return
	}
}
