package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/dashboard"
	"log"
	"net/http"
	"strconv"
)

func (h Handlers) ToggleCompleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("toggle complete: endpoint hit without HTMX")
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("toggle complete: no database connection provided to handler")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Println("toggle complete: no id")
		return
	}

	var (
		t, newTask   task.Task
		err error
	)
	t.ID, err = strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		log.Printf("toggle complete: %s\n", err)
		return
	}

	t, newTask, err = sqlite.ToggleTaskComplete(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		log.Printf("toggle complete: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = dashboard.ListItem(t, newTask).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("toggle complete: %s\n", err)
		return
	}
}
