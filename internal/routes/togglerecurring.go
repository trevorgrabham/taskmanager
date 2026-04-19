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

// Renders or removes the RecurringPeriodInput template.
func (h Handlers) ToggleRecurringTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("toggle recurring: endpoint hit without HTMX")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parsing error", http.StatusBadRequest)
		log.Println("toggle recurring: error parsing form")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if r.FormValue("recurring") == "on" {
		idString := r.URL.Query().Get("id")
		if idString == "" {
			http.Error(w, "Error bad id", http.StatusBadRequest)
			log.Println("toggle recurring: no id")
			return
		}

		var t task.Task
		t.ID, err = strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Error bad id", http.StatusBadRequest)
			log.Printf("toggle recurring: %s\n", err)
			return
		}

		if t.ID != 0 {
			t, err = sqlite.TaskByID(h.DB, t.ID)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				log.Printf("toggle recurring: %s\n", err)
				return
			}
		}

		err = shared.RecurringPeriodInput(t).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			log.Printf("toggle-recurring: %s\n", err)
			return
		}
	} else {
		fmt.Fprint(w, "")
	}
}
