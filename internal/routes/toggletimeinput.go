package routes

import (
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/addtaskform"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h Handlers) ToggleTimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("toggle recurring: endpoint hit without HTMX")
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("nil database")
		return
	}

	var err error
	switch strings.TrimPrefix(r.URL.Path, "/toggle-time/") {
	case "add-task":
		w.Header().Set("Content-Type", "text/html")
		if r.URL.Query().Get("all-day") == "on" {
			fmt.Fprint(w, "")
		} else {
			err = addtaskform.TimeInput(task.Task{}).Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Rendering error", http.StatusInternalServerError)
				log.Printf("toggle time: %s\n", err)
				return
			}
		}
	case "edit-task":
		if r.URL.Query().Get("all-day") == "on" {
			fmt.Fprint(w, "")
		} else {
			idString := r.URL.Query().Get("id")
			if idString == "" {
				http.Error(w, "Error no id", http.StatusBadRequest)
				log.Println("toggle time: no id")
				return
			}

			var (
				id int
				t  task.Task
			)
			id, err = strconv.Atoi(idString)
			if err != nil {
				http.Error(w, "Error bad id", http.StatusBadRequest)
				log.Printf("toggle time: %s\n", err)
				return
			}

			t, err = sqlite.TaskByID(h.DB, id)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				log.Printf("toggle time: %s\n", err)
				return
			}

			err = addtaskform.TimeInput(t).Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Rendering error", http.StatusInternalServerError)
				log.Printf("toggle time: %s\n", err)
				return
			}
		}
	default:
		http.Error(w, "Error bad path", http.StatusBadRequest)
		log.Printf("toggle time: unknown path %s\n", r.URL.RawPath)
		return
	}
}
