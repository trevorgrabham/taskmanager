package routes

import (
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views"
	"log"
	"net/http"
	"time"
)

func (h Handlers) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	dueDateString := r.URL.Query().Get("due-date")
	if dueDateString != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateString)
		if err != nil {
			http.Error(w, "Bad due date parameter", http.StatusBadRequest)
			log.Printf("add task form: bad due date parameter %s", dueDateString)
			return
		}
		t.DueDate = task.TaskDueDate(dueDate)
	}
	w.Header().Set("Content-Type", "text/html")
	err := views.AddTaskForm(t).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("index: %s\n", err)
		return
	}
}
