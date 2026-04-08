package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("delete task: nil database")
		return
	}

	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) != 3 {
		http.Error(w, "Bad path", http.StatusBadRequest)
		log.Printf("delete task: bad path %s\n", r.URL.Path)
		return
	}

	id, err := strconv.Atoi(splitPath[2])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		log.Printf("delete task: bad id %s\n", err)
		return
	}

	err = sqlite.DeleteTask(h.DB, task.Task{ID: id})
	if err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		log.Printf("delete task: %s\n", err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
