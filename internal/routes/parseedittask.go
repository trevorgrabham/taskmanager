package routes

import (
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/taskinfo"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h Handlers) ParseEditTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("parsing edit task: endpoint hit without HTMX")
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("parsing edit task: no database provided to handler")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse New Task Form", http.StatusBadRequest)
		log.Println("parsing edit task: error parsing form")
		return
	}

	var (
		t       task.Task
		dueDate time.Time
	)
	t.Title = r.FormValue("title")
	if t.Title == "" {
		http.Error(w, "Error no title", http.StatusBadRequest)
		log.Println("parsing edit task: no title provided")
		return
	}

	t.Category = r.FormValue("category")
	t.Description = r.FormValue("description")
	dueDateString := r.FormValue("due-date")
	if dueDateString != "" {
		dueDate, err = time.ParseInLocation("2006-01-02 3:04PM", dueDateString+" 11:59PM", time.Local)
		if err != nil {
			http.Error(w, "Error bad due date", http.StatusBadRequest)
			log.Printf("parsing edit task: %s\n", err)
			return
		}

		t.DueDate = task.TaskDueDate(dueDate)
	}

	idString := r.FormValue("id")
	if idString == "" {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Println("parsing edit task: no id")
		return
	}

	t.ID, err = strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("parsing edit task: %s\n", err)
		return
	}

	recurringPeriodValue := r.FormValue("period-value")
	recurringPeriodUnit := r.FormValue("period-unit")
	if recurringPeriodValue != "" && recurringPeriodUnit != "" {
		if recurringPeriodUnit != "days" && recurringPeriodUnit != "weeks" && recurringPeriodUnit != "months" {
			http.Error(w, "Bad value for recurring unit", http.StatusBadRequest)
			log.Printf("parsing edit task: bad value for recurring period unit %s\n", recurringPeriodUnit)
			return
		}

		var value int
		value, err = strconv.Atoi(recurringPeriodValue)
		if err != nil {
			http.Error(w, "Bad value for recurring value", http.StatusBadRequest)
			log.Printf("parsing edit task: bad value for recurring period value %s\n", recurringPeriodValue)
			return
		}

		t.RecurringPeriod = fmt.Sprintf("%d %s", value, recurringPeriodUnit)
	}

	err = sqlite.UpdateTask(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("parsing edit task: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = taskinfo.TaskInfo(t).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("parsing edit task: %s\n", err)
		return
	}
}
