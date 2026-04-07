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

func (h Handlers) ParseAddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Println("ParseAddTaskHandler(): endpoint hit without HTMX")
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("ParseAddTaskHandler(): no database provided to handler")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse New Task Form", http.StatusBadRequest)
		log.Println("ParseAddTaskHandler(): error parsing form")
		return
	}

	var (
		t       task.Task
		dueDate time.Time
	)
	t.Title = r.FormValue("title")
	if t.Title == "" {
		http.Error(w, "No 'title' provided", http.StatusBadRequest)
		log.Println("ParseAddTaskHandler(): no title provided")
		return
	}

	t.Category = r.FormValue("category")
	t.Description = r.FormValue("description")
	dueDateString := r.FormValue("due-date")
	if dueDateString != "" {
		dueDate, err = time.ParseInLocation("2006-01-02 3:04PM", dueDateString+" 11:59PM", time.Local)
		if err != nil {
			http.Error(w, "Error parsing due date", http.StatusBadRequest)
			log.Printf("ParseAddTaskHandler(): error parsing due date")
			return
		}

		t.DueDate = task.TaskDueDate(dueDate)
	}

	var taskID int
	taskID, err = sqlite.AddTask(h.DB, t)
	if err != nil {
		http.Error(w, "Error adding task", http.StatusBadRequest)
		log.Printf("ParseAddTaskHandler(): error adding task to DB")
		return
	}
	log.Println(t.Title)
	log.Println(t.Category)
	log.Println(t.Description)
	log.Println(t.DueDate.String())

	recurringPeriodValue := r.FormValue("period-value")
	recurringPeriodUnit := r.FormValue("period-unit")
	if recurringPeriodValue != "" && recurringPeriodUnit != "" {
		if recurringPeriodUnit != "days" && recurringPeriodUnit != "weeks" && recurringPeriodUnit != "months" {
			http.Error(w, "Bad value for recurring unit", http.StatusBadRequest)
			log.Printf("ParseAddTaskHandler(): bad value for recurring period unit %s", recurringPeriodUnit)
			return
		}

		var value int
		value, err = strconv.Atoi(recurringPeriodValue)
		if err != nil {
			http.Error(w, "Bad value for recurring value", http.StatusBadRequest)
			log.Printf("ParseAddTaskHandler(): bad value for recurring period value %s", recurringPeriodValue)
			return
		}

		err = sqlite.AddRecurringTask(h.DB, task.RecurringTask{TaskID: taskID, Period: fmt.Sprintf("%d %s", value, recurringPeriodUnit)})
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Printf("ParseAddTaskHandler(): %s", err)
			return
		}
	}

	w.Header().Set("HX-Redirect", "/")
}
