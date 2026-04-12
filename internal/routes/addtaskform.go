package routes

import (
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/addtaskform"
	sqlite "local/taskmanager/internal/db"
	"log"
	"net/http"
	"time"
)

// GET: Renders the add-task form
// POST: Parses the add-task form
func (h Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { 
		h.ParseAddTask(w, r)
		return
	}

	var t task.Task
	dueDateString := r.URL.Query().Get("due-date")
	if dueDateString != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateString)
		if err != nil {
			http.Error(w, "Bad due date parameter", http.StatusBadRequest)
			log.Printf("add task form: bad due date parameter %s\n", dueDateString)
			return
		}
		t.DueDate = task.TaskDueDate(dueDate)
	}
	t.Category = r.URL.Query().Get("category")

	userCategories, err := sqlite.Categories(h.DB)
	if err != nil { 
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("add task form: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = addtaskform.AddTaskForm(t, userCategories).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("add task form: %s\n", err)
		return
	}
}
