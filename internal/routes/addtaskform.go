package routes

import (
	"fmt"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
)

// GET: Renders the add-task form
func (h Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok bool
		t task.Task
		userCategories []string
		dueDateString = fmt.Sprintf("%s 00:00", r.URL.Query().Get("due-date"))
	)
	if ok = h.checkConnection(w); !ok {
		return
	}
	
	if t.DueDate, ok = h.parseDueDate(w, dueDateString); !ok { 
		return 
	}

	t.Category = r.URL.Query().Get("category")

	if userCategories, ok = h.getUserCategorySuggestions(w); !ok {
		return 
	}

	w.Header().Set("Content-Type", "text/html")
	err := shared.TaskForm(t, userCategories).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("add task handler: %s\n", err)
		return
	}
}
