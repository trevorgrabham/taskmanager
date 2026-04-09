package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
	"strconv"
)

// Renders task-info page
func (h Handlers) TaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("task info: nil database")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Println("task info: no id")
		return
	}

	var (
		t   task.Task
		err error
	)
	t.ID, err = strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("task info: %s\n", err)
		return
	}

	t, err = sqlite.TaskByID(h.DB, t.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("task info: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = taskinfo.TaskInfo(t).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("task info: %s", err)
		return
	}
}
