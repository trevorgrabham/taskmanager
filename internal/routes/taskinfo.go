package routes

import (
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
)

// Renders task-info page
func (h Handlers) TaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok bool
		t  task.Task
	)
	if ok = h.checkConnection(w); !ok {
		return
	}

	if t.ID, ok = h.parseID(w, r.URL.Query().Get("id")); !ok {
		return
	}

	if t, ok = h.getTask(w, t.ID); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := taskinfo.TaskInfo(t).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("task info: %s", err)
		return
	}
}
