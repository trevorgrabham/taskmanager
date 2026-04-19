package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
	"time"
)

// GET: Renders the add-task form
func (h Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("add task handler: no database connection")
		return
	}

	var t task.Task
	dueDateString := r.URL.Query().Get("due-date")
	if dueDateString != "" {
		dueDate, err := time.ParseInLocation("2006-01-02", dueDateString, time.Local)
		if err != nil {
			http.Error(w, "Error bad due date", http.StatusBadRequest)
			log.Printf("add task handler: bad due date parameter %s\n", dueDateString)
			return
		}
		t.DueDate = task.TaskDueDate(dueDate)
	}
	t.Category = r.URL.Query().Get("category")

	userCategories, err := sqlite.Categories(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("add task handler: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = shared.TaskForm(t, userCategories).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("add task handler: %s\n", err)
		return
	}
}
