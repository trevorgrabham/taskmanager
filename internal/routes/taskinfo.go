package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
	"strconv"
)

func (h Handlers) TaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("delete task: nil database")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Bad route", http.StatusBadRequest)
		log.Printf("task info: bad path %s\n", r.URL.RawPath)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Bad route", http.StatusBadRequest)
		log.Printf("task info: bad id %s\n", idString)
		return
	}

	var t task.Task
	t, err = sqlite.TaskByID(h.DB, id)
	if err != nil {
		http.Error(w, "Error getting task info", http.StatusInternalServerError)
		log.Printf("task info: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	taskinfo.TaskInfo(t).Render(r.Context(), w)
}
