package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"log"
	"net/http"
	"strconv"
)

// Hits the backend db.DeleteTask() and redirects to "/"
func (h Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("delete task: nil database")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Println("delete task: no id")
		return
	}

	var (
		t   task.Task
		err error
	)
	t.ID, err = strconv.Atoi(idString)
	if idString == "" {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("delete task: %s\n", err)
		return
	}

	err = sqlite.DeleteTask(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		log.Printf("delete task: %s\n", err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
