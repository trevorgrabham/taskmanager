package routes

import (
	"fmt"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
)

func (h Handlers) ToggleTimeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok bool
		t  task.Task
	)
	if ok = h.checkHXRequest(w, r); !ok {
		return
	}

	if ok := h.checkConnection(w); !ok {
		return
	}

	if r.URL.Query().Get("all-day") == "on" {
		fmt.Fprint(w, "")
	} else {
		if t.ID, ok = h.parseID(w, r.URL.Query().Get("id")); !ok {
			return
		}

		if t.ID != 0 {
			if t, ok = h.getTask(w, t.ID); !ok {
				return
			}
		}

		err := shared.TimeInput(t).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			log.Printf("toggle time handler: %s\n", err)
			return
		}
	}
}
