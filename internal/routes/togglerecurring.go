package routes

import (
	"fmt"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
)

// Renders or removes the RecurringPeriodInput template.
func (h Handlers) ToggleRecurringTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok  bool
		t   task.Task
		err error
	)
	if ok = h.checkHXRequest(w, r); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if r.URL.Query().Get("recurring") == "on" {
		if t.ID, ok = h.parseID(w, r.URL.Query().Get("id")); !ok {
			return
		}

		if t.ID != 0 {
			if t, ok = h.getTask(w, t.ID); !ok {
				return
			}
		}

		err = shared.RecurringPeriodInput(t).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			log.Printf("toggle-recurring: %s\n", err)
			return
		}
	} else {
		fmt.Fprint(w, "")
	}
}
