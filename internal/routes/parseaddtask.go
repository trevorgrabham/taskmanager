package routes

import (
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Parses the add-task POST request and adds it to our backend. Redirects to "/"
func (h Handlers) ParseAddTask(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("parsing add task: no database provided to handler")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse New Task Form", http.StatusBadRequest)
		log.Printf("parsing add task: %s\n", err)
		return
	}

	var (
		t       task.Task
		dueDate time.Time
	)
	t.Title = r.FormValue("title")
	if t.Title == "" {
		http.Error(w, "Error no title", http.StatusBadRequest)
		log.Println("parsing add task: no title provided")
		return
	}

	t.Category = r.FormValue("category")
	t.Description = r.FormValue("description")
	dueDateString := r.FormValue("due-date")
	if dueDateString != "" {
		dueDate, err = time.ParseInLocation("2006-01-02 3:04PM", dueDateString+" 11:59PM", time.Local)
		if err != nil {
			http.Error(w, "Error parsing due date", http.StatusBadRequest)
			log.Printf("parsing add task: %s\n", err)
			return
		}

		t.DueDate = task.TaskDueDate(dueDate)
	}

	recurringPeriodValue := r.FormValue("period-value")
	recurringPeriodUnit := r.FormValue("period-unit")
	if recurringPeriodValue != "" && recurringPeriodUnit != "" {
		if recurringPeriodUnit != "days" && recurringPeriodUnit != "weeks" && recurringPeriodUnit != "months" {
			http.Error(w, "Bad value for recurring unit", http.StatusBadRequest)
			log.Printf("parsing add task: bad value for recurring period unit %s\n", recurringPeriodUnit)
			return
		}

		var value int
		value, err = strconv.Atoi(recurringPeriodValue)
		if err != nil {
			http.Error(w, "Bad value for recurring value", http.StatusBadRequest)
			log.Printf("parsing add task: bad value for recurring period value %s\n", recurringPeriodValue)
			return
		}

		t.RecurringPeriod = fmt.Sprintf("%d %s", value, recurringPeriodUnit)
	}

	_, err = sqlite.AddTask(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("parsing add task: %s\n", err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}
