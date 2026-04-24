package routes

import (
	"errors"
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
	"strings"
)

// Renders or removes the RecurringPeriodInput template.
func (h Handlers) ToggleRecurringTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller                     = "ToggleRecurringTaskHandler"
		err                        error
		sessionID, recurringToggle string
		user                       services.User
		t                          services.Task
		split                      []string
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
	if recurringToggle = r.URL.Query().Get("recurring"); recurringToggle != "on" {
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

	if split = strings.Split(t.RecurringPeriod, " "); len(split) != 2 {
		// TODO: wrap the returned error
		h.HandleError(w, errors.New("badly formed recurring period"))
		return
	}
	err = views.RecurringPeriodInput(split[0], split[1]).Render(r.Context(), w)
	if err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "RecurringPeriodInput", err))
		return
	}
}
