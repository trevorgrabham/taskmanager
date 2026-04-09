package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
	"strconv"
)

// POST request: Parses form data and calls routes/ParseEditTaskHandler() to update the backend
// GET request: Renders task-info page in view mode
func (h Handlers) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("delete task: nil database")
		return
	}

	if r.Method == http.MethodPost {
		h.ParseEditTask(w, r)
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Println("task info: no id")
		return
	}

	var (
		t   task.Task
		err error
	)
	t.ID, err = strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Bad route", http.StatusBadRequest)
		log.Printf("task info: bad id %s\n", idString)
		return
	}

	t, err = sqlite.TaskByID(h.DB, t.ID)
	if err != nil {
		http.Error(w, "Error getting task info", http.StatusInternalServerError)
		log.Printf("task info: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = taskinfo.EditTaskInfo(t).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("task info: %s", err)
		return
	}
}
