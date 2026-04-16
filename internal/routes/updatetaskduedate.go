package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/dashboard"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h Handlers) UpdateTaskDueDateHandler(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("update task duedate: nil database")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Println("update task duedate: no id")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Error no id", http.StatusBadRequest)
		log.Println("update task duedate: no id")
		return
	}
	newDate := r.URL.Query().Get("date")
	if newDate == "" {
		err = sqlite.UpdateDueDate(h.DB, task.Task{ID: id})
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Printf("update task duedate: %s\n", err)
			return
		}

		var unscheduledTasks task.TaskList
		unscheduledTasks, err = sqlite.UnscheduledTasks(h.DB)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Printf("update task duedate: %s\n", err)
			return
		}

		unscheduledMap := make(map[string]task.TaskList)
		for _, t := range unscheduledTasks {
			unscheduledMap[t.Category] = append(unscheduledMap[t.Category], t)
		}
		err = dashboard.UnscheduledTasks(unscheduledMap).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			log.Printf("update task duedate: %s\n", err)
			return
		}

		return
	}

	var newDueDate time.Time
	newDueDate, err = time.ParseInLocation("2006-01-02", newDate, time.Local)
	if err != nil {
		http.Error(w, "Error bad date", http.StatusBadRequest)
		log.Printf("update task duedate: %s\n", err)
		return
	}

	err = sqlite.UpdateDueDate(h.DB, task.Task{ID: id, DueDate: task.TaskDueDate(newDueDate)})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("update task duedate: %s\n", err)
		return
	}

	start := time.Date(newDueDate.Year(), newDueDate.Month(), newDueDate.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(newDueDate.Year(), newDueDate.Month(), newDueDate.Day(), 23, 59, 0, 0, time.Local)
	var daysTasks task.TaskList
	daysTasks, err = sqlite.QueryTasks(h.DB, sqlite.TaskQueryParams{WhichTasks: sqlite.IncTasks, From: start, To: end})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("update task duedate: %s\n", err)
		return
	}

	dayMap := make(map[string]task.TaskList)
	for _, t := range daysTasks {
		dayMap[t.Category] = append(dayMap[t.Category], t)
	}
	err = dashboard.DayTaskList(dayMap, newDueDate).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("update task duedate: %s\n", err)
		return
	}
}
