package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
	"strconv"
)

// GET request: Renders task-info page in view mode
func (h Handlers) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("edit task handler: no database connection")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Println("edit task handler: no id")
		return
	}

	var (
		t              task.Task
		userCategories []string
		err            error
	)
	t.ID, err = strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("edit task handler: bad id %s\n", idString)
		return
	}

	t, err = sqlite.TaskByID(h.DB, t.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("edit task handler: %s\n", err)
		return
	}

	userCategories, err = sqlite.Categories(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("edit task handler: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = shared.TaskForm(t, userCategories).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("edit task handler: %s\n", err)
		return
	}
}
