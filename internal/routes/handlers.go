package routes

import (
	"database/sql"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Handlers struct {
	DB *sql.DB
}

func (h Handlers) addTask(w http.ResponseWriter, t task.Task) (ok bool) {
	_, err := sqlite.AddTask(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return false
	}

	return true
}

func (h Handlers) checkConnection(w http.ResponseWriter) (ok bool) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: no database connection\n", getCallingFunc(2))
		return false
	}
	return true
}

func (h Handlers) checkHXRequest(w http.ResponseWriter, r *http.Request) (ok bool) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Endpoint expected an HTMX request", http.StatusBadRequest)
		log.Printf("%s: endpoint hit without HTMX\n", getCallingFunc(2))
		return false
	}

	return true
}

func (h Handlers) checkKnownPath(w http.ResponseWriter, r *http.Request) (ok bool) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		log.Printf("%s: unknown path %s\n", getCallingFunc(2), r.URL.Path)
		return false
	}

	return true
}

func (h Handlers) checkMethod(w http.ResponseWriter, r *http.Request, method string) (ok bool) {
	if r.Method != method {
		http.Error(w, fmt.Sprintf("Expected %s request", method), http.StatusBadRequest)
		log.Printf("%s: not a %s request", getCallingFunc(2), method)
		return false
	}

	return true
}

func (h Handlers) completeTask(w http.ResponseWriter, id int) (ok bool) {
	err := sqlite.CompleteTask(h.DB, task.Task{ID: id})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return false
	}

	return true
}

func (h Handlers) deleteTask(w http.ResponseWriter, id int) (ok bool) {
	err := sqlite.DeleteTask(h.DB, task.Task{ID: id})
	if err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return false
	}

	return true
}

func (h Handlers) getFavCategory(weekOfTasks map[int]map[string]task.TaskList) task.TaskList {
	var favCategory string
	tasksByCategory := make(map[string]task.TaskList)
	for dateKey := range weekOfTasks {
		for catKey := range weekOfTasks[dateKey] {
			for _, t := range weekOfTasks[dateKey][catKey] {
				tasksByCategory[t.Category] = append(tasksByCategory[t.Category], t)
			}
		}
	}

	for k := range tasksByCategory {
		if favCategory == "" {
			favCategory = k
			continue
		}

		if len(tasksByCategory[k]) > len(tasksByCategory[favCategory]) {
			favCategory = k
		}
	}
	return tasksByCategory[favCategory]
}

func (h Handlers) getOverdueTasks(w http.ResponseWriter) (overdueTasks map[string]task.TaskList, ok bool) {
	overdue, err := sqlite.OverdueTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s", getCallingFunc(2), err)
		return nil, false
	}

	overdueTasks = make(map[string]task.TaskList)
	for _, t := range overdue {
		overdueTasks[t.Category] = append(overdueTasks[t.Category], t)
	}

	return overdueTasks, true
}

func (h Handlers) getTask(w http.ResponseWriter, id int) (t task.Task, ok bool) {
	var err error
	t, err = sqlite.TaskByID(h.DB, id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return t, false
	}

	return t, true
}

func (h Handlers) getTasksForDay(w http.ResponseWriter, date time.Time) (daysTasks map[string]task.TaskList, ok bool) {
	tasks, err := sqlite.ListDailyTasks(h.DB, date)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return nil, false
	}

	daysTasks = make(map[string]task.TaskList)
	for _, t := range tasks {
		daysTasks[t.Category] = append(daysTasks[t.Category], t)
	}

	return daysTasks, true
}

func (h Handlers) getUnscheduledTasks(w http.ResponseWriter) (unscheduledTasks map[string]task.TaskList, ok bool) {
	unscheduled, err := sqlite.UnscheduledTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s", getCallingFunc(2), err)
		return nil, false
	}

	unscheduledTasks = make(map[string]task.TaskList)
	for _, t := range unscheduled {
		unscheduledTasks[t.Category] = append(unscheduledTasks[t.Category], t)
	}

	return unscheduledTasks, true
}

func (h Handlers) getUpcomingWeek(w http.ResponseWriter) (weekOfTasks map[int]map[string]task.TaskList, ok bool) {
	var err error
	weekOfTasks, err = sqlite.ListWeekOfTasks(h.DB, time.Now())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s", getCallingFunc(2), err)
		return nil, false
	}

	return weekOfTasks, true
}

func (h Handlers) getUserCategorySuggestions(w http.ResponseWriter) (categorySuggestions []string, ok bool) {
	var err error
	categorySuggestions, err = sqlite.Categories(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return nil, false
	}

	return categorySuggestions, true
}

func (h Handlers) parseDueDate(w http.ResponseWriter, dateString string) (dueDate task.TaskDueDate, ok bool) {
	if dateString != "" {
		date, err := time.ParseInLocation("2006-01-02 15:04", dateString, time.Local)
		if err != nil {
			http.Error(w, "Error bad due date", http.StatusBadRequest)
			log.Printf("%s: bad due date parameter %s\n", getCallingFunc(2), dateString)
			return dueDate, false
		}
		dueDate = task.TaskDueDate(date)
	}
	return dueDate, true
}

func (h Handlers) parseID(w http.ResponseWriter, idString string) (id int, ok bool) {
	if idString == "" {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Printf("%s: no id", getCallingFunc(2))
		return 0, false
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error bad id", http.StatusBadRequest)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return 0, false
	}

	return id, true
}

func (h Handlers) parseRecurringPeriod(w http.ResponseWriter, valueString, unitString string) (period string, ok bool) {
	if valueString == "" || unitString == "" {
		return "", true
	}

	if unitString != "days" && unitString != "weeks" && unitString != "months" {
		http.Error(w, "Bad value for recurring unit", http.StatusBadRequest)
		log.Printf("%s: bad value for recurring period unit %s\n", getCallingFunc(2), unitString)
		return "", false
	}

	value, err := strconv.Atoi(valueString)
	if err != nil {
		http.Error(w, "Bad value for recurring value", http.StatusBadRequest)
		log.Printf("%s: bad value for recurring period value %s\n", getCallingFunc(2), valueString)
		return "", false
	}

	return fmt.Sprintf("%d %s", value, unitString), true
}

func (h Handlers) parseTitle(w http.ResponseWriter, r *http.Request) (title string, ok bool) {
	title = r.FormValue("title")
	if title == "" {
		http.Error(w, "Error no title", http.StatusBadRequest)
		log.Printf("%s: no title provided", getCallingFunc(2))
		return "", false
	}

	return title, true
}

func (h Handlers) toggleTaskComplete(w http.ResponseWriter, t task.Task) (toggledTask, nextTask task.Task, ok bool) {
	var err error
	toggledTask, nextTask, err = sqlite.ToggleTaskComplete(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusBadRequest)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return task.Task{}, task.Task{}, false
	}

	return toggledTask, nextTask, true
}

func (h Handlers) updateDueDate(w http.ResponseWriter, t task.Task) (ok bool) {
	err := sqlite.UpdateDueDate(h.DB, t)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return false
	}

	return true
}

func (h Handlers) updateTask(w http.ResponseWriter, t task.Task) (ok bool) {
	if err := sqlite.UpdateTask(h.DB, t); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("%s: %s\n", getCallingFunc(2), err)
		return false
	}

	return true
}

func getCallingFunc(skip int) string {
	if skip < 0 {
		return ""
	}

	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}

	split := strings.Split(fn.Name(), ".")
	return split[len(split)-1]
}
