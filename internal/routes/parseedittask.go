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

func (h Handlers) ParseEditTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("ParseEditTaskHandler(): endpoint hit without HTMX")
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("ParseEditTaskHandler(): no database provided to handler")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse New Task Form", http.StatusBadRequest)
		log.Println("ParseEditTaskHandler(): error parsing form")
		return
	}

	var (
		t       task.Task
		dueDate time.Time
	)
	t.Title = r.FormValue("title")
	if t.Title == "" {
		http.Error(w, "No 'title' provided", http.StatusBadRequest)
		log.Println("ParseEditTaskHandler(): no title provided")
		return
	}

	t.Category = r.FormValue("category")
	t.Description = r.FormValue("description")
	dueDateString := r.FormValue("due-date")
	if dueDateString != "" {
		dueDate, err = time.ParseInLocation("2006-01-02 3:04PM", dueDateString+" 11:59PM", time.Local)
		if err != nil {
			http.Error(w, "Error parsing due date", http.StatusBadRequest)
			log.Printf("ParseEditTaskHandler(): error parsing due date")
			return
		}

		t.DueDate = task.TaskDueDate(dueDate)
	}

	idString := r.FormValue("id")
	if idString == "" {
		http.Error(w, "Error parsing id", http.StatusBadRequest)
		log.Printf("ParseEditTaskHandler(): error parsing id")
		return
	}

	var id int
	id, err = strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error parsing id", http.StatusBadRequest)
		log.Printf("ParseEditTaskHandler(): error parsing id")
		return
	}

	t.ID = id
	recurringPeriodValue := r.FormValue("period-value")
	recurringPeriodUnit := r.FormValue("period-unit")
	if recurringPeriodValue != "" && recurringPeriodUnit != "" {
		if recurringPeriodUnit != "days" && recurringPeriodUnit != "weeks" && recurringPeriodUnit != "months" {
			http.Error(w, "Bad value for recurring unit", http.StatusBadRequest)
			log.Printf("ParseEditTaskHandler(): bad value for recurring period unit %s", recurringPeriodUnit)
			return
		}

		var value int
		value, err = strconv.Atoi(recurringPeriodValue)
		if err != nil {
			http.Error(w, "Bad value for recurring value", http.StatusBadRequest)
			log.Printf("ParseEditTaskHandler(): bad value for recurring period value %s", recurringPeriodValue)
			return
		}

		t.RecurringPeriod = fmt.Sprintf("%d %s", value, recurringPeriodUnit)
	}

	fmt.Println(t.Title, t.Category, t.Description, t.DueDate, t.RecurringPeriod)
	err = sqlite.UpdateTask(h.DB, t)
	if err != nil {
		http.Error(w, "Error adding task", http.StatusBadRequest)
		log.Printf("ParseEditTaskHandler(): error adding task to DB")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	taskinfo.TaskInfo(t).Render(r.Context(), w)
}
