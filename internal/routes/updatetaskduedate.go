package routes

import (
	"fmt"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views/dashboard"
	"log"
	"net/http"
	"time"
)

func (h Handlers) UpdateTaskDueDateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ok               bool
		t                task.Task
		unscheduledTasks map[string]task.TaskList
		dueDate          time.Time
		daysTasks        map[string]task.TaskList
	)
	if ok = h.checkConnection(w); !ok {
		return
	}

	if t.ID, ok = h.parseID(w, r.URL.Query().Get("id")); !ok {
		return
	}
	newDate := r.URL.Query().Get("duedate")
	if newDate == "" {
		if ok = h.updateDueDate(w, t); !ok {
			return
		}

		if unscheduledTasks, ok = h.getUnscheduledTasks(w); !ok {
			return
		}

		err := dashboard.UnscheduledTasks(unscheduledTasks).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			log.Printf("update task duedate: %s\n", err)
			return
		}

		return
	}

	if t, ok = h.getTask(w, t.ID); !ok {
		return
	}

	dueDate = time.Time(t.DueDate)
	if t.DueDate, ok = h.parseDueDate(w, fmt.Sprintf("%s %2d:%2d", newDate, dueDate.Hour(), dueDate.Minute())); !ok {
		return
	}

	if ok = h.updateDueDate(w, t); !ok {
		return
	}

	dueDate = time.Time(t.DueDate)
	if daysTasks, ok = h.getTasksForDay(w, dueDate); !ok {
		return
	}

	err := dashboard.DayTaskList(daysTasks, dueDate).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Rendering error", http.StatusInternalServerError)
		log.Printf("update task duedate: %s\n", err)
		return
	}
}
