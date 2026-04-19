package routes

import (
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
)

// GET request: Renders task-info page in view mode
func (h Handlers) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		categories []string
		id         int
		ok         bool
		t          task.Task
	)

	if ok = h.checkConnection(w); !ok {
		return
	}

	if id, ok = h.parseID(w, r.URL.Query().Get("id")); !ok {
		return
	}

	if t, ok = h.getTask(w, id); !ok {
		return
	}

	if categories, ok = h.getUserCategorySuggestions(w); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := shared.TaskForm(t, categories).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("edit task handler: %s\n", err)
		return
	}
}
