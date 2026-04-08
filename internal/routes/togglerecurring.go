package routes

import (
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/addtaskform"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h Handlers) ToggleRecurringTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("toggle recurring: endpoint hit without HTMX")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse New Task Form", http.StatusBadRequest)
		log.Println("toggle recurring: error parsing form")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if r.FormValue("recurring") == "on" {
		switch strings.TrimPrefix(r.URL.Path, "/toggle-recurring/") {
		case "add-task":
			err = addtaskform.RecurringPeriodInput().Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Rendering error", http.StatusInternalServerError)
				log.Printf("toggle-recurring/add-task: %s\n", err)
				return
			}
		case "edit-task":
			idString := r.URL.Query().Get("id")
			if idString == "" {
				http.Error(w, "Bad path", http.StatusBadRequest)
				log.Println("toggle recurring: no id")
				return
			}

			var id int
			id, err = strconv.Atoi(idString)
			if err != nil {
				http.Error(w, "Bad path", http.StatusBadRequest)
				log.Printf("toggle recurring: bad path %s\n", r.URL.RawPath)
				return
			}

			var t task.Task
			t, err = sqlite.TaskByID(h.DB, id)
			if err != nil {
				http.Error(w, "Bad path", http.StatusBadRequest)
				log.Printf("toggle recurring: bad path %s\n", r.URL.RawPath)
				return
			}

			err = taskinfo.RecurringPeriodInput(t).Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Rendering error", http.StatusInternalServerError)
				log.Printf("toggle-recurring/edit-task: %s\n", err)
				return
			}
		default:
			http.Error(w, "Bad path", http.StatusBadRequest)
			log.Printf("toggle recurring: bad path %s\n", r.URL.Path)
			return
		}
	} else {
		fmt.Fprint(w, "")
	}
}
