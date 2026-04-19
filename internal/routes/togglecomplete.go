package routes

import (
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/dashboard"
	"log"
	"net/http"
)

func (h Handlers) ToggleCompleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok bool
		t, newTask  task.Task
	)
	if ok = h.checkHXRequest(w, r); !ok {
		return
	}

	if ok = h.checkConnection(w); !ok {
		return
	}

	if t.ID, ok = h.parseID(w, r.URL.Query().Get("id")); !ok { return }

	if t, newTask, ok = h.toggleTaskComplete(w, t); !ok { return }

	w.Header().Set("Content-Type", "text/html")
	err := dashboard.ListItem(t, newTask).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("toggle complete: %s\n", err)
		return
	}
}
