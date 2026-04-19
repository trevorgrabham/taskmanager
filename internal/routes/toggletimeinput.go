package routes

import (
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
	"strconv"
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
	if r.URL.Query().Get("all-day") == "on" {
		fmt.Fprint(w, "")
	} else {
		idString := r.URL.Query().Get("id")
		if idString == "" {
			http.Error(w, "Error no id", http.StatusBadRequest)
			log.Println("toggle time handler: no id")
			return
		}

		var (
			id int
			t  task.Task
		)
		id, err = strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Error bad id", http.StatusBadRequest)
			log.Printf("toggle time handler: %s\n", err)
			return
		}

		if id != 0 {
			t, err = sqlite.TaskByID(h.DB, id)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				log.Printf("toggle time handler: %s\n", err)
				return
			}
		}

		err = shared.TimeInput(t).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			log.Printf("toggle time handler: %s\n", err)
			return
		}
	}
}
