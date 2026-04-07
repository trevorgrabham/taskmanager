package routes

import (
	"fmt"
	"local/taskmanager/internal/views/addtaskform"
	"log"
	"net/http"
)

func (h Handlers) ToggleRecurringTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("AddTaskHandler(): endpoint hit without HTMX")
		return
	}

	log.Println("/recurring-toggle endpoint hit")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse New Task Form", http.StatusBadRequest)
		log.Println("AddTaskHandler(): error parsing form")
		return
	}

	log.Println(r.FormValue("recurring"))

	w.Header().Set("Content-Type", "text/html")
	if r.FormValue("recurring") == "on" {
		fmt.Fprint(w, addtaskform.RecurringPeriodInput().Render(r.Context(), w))
	} else {
		fmt.Fprint(w, "")
	}
}
