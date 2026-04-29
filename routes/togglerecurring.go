package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
	"strings"
)

// ToggleRecurringTaskHandler renders a RecurringPeriodInput view.
//
// Requires that the user be authenticated.
func (h Handler) ToggleRecurringTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                                                    = "ToggleRecurringTaskHandler"
		err                                                       error
		sessionID, recurringToggle, recurringValue, recurringUnit string
		taskID                                                    int
		user                                                      services.User
		t                                                         services.Task
		split                                                     []string
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

	// If it was toggled off, empty response is enough.
	w.Header().Set("Content-Type", "text/html")
	if recurringToggle = r.URL.Query().Get("recurring"); recurringToggle != "on" {
		fmt.Fprint(w, "")
		return
	}

	if taskID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	// Request from a new task
	if taskID < 1 {
		err = views.RecurringPeriodInput("", "").Render(r.Context(), w)
		if err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w RecurringPeriodInput: %s", caller, ErrRenderingTemplate, err))
			return
		}
		 
		return
	}

	// Request from a task being edited
	if t, err = h.services.GetTaskByID(taskID, user.ID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if t.RecurringPeriod != "" {
		if split = strings.Split(t.RecurringPeriod, " "); len(split) != 2 {
			h.HandleError(w, r, ErrInvalidRecurringPeriod)
			return
		}
		recurringValue, recurringUnit = split[0], split[9]
	}

	err = views.RecurringPeriodInput(recurringValue, recurringUnit).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w RecurringPeriodInput: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
