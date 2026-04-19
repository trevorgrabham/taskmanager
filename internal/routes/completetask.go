package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"log"
	"net/http"
	"strconv"
)

func (h Handlers) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("completing task: no database set up")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Println("completing task: no id")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("completing task: %s\n", err)
		return
	}

	err = sqlite.CompleteTask(h.DB, task.Task{ID: id})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("completing task: %s\n", err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
