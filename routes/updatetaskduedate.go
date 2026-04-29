package routes

import (
	"fmt"
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"net/http"
	"time"
)

// UpdateTaskDueDateHandler updates the DueDate for the task identified by taskID to the due-date provided. Renders the list of tasks on the new due-date, or unscheduled tasks if the due-date is empty.
//
// If due-date is non-empty and not a valid date-time, no update is performed and a non-2XX status code is returned.
// If due-date is empty, the DueDate is removed from the Task.
//
// Requires that the user be authenticated.
func (h Handler) UpdateTaskDueDateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller               = "UpdateTaskDueDateHandler"
		err                  error
		sessionID, dayString string
		newDueDay            time.Time
		user                 services.User
		updatedTask          services.Task
		tasks                services.TaskList
		taskID               int
		unscheduledData      views.UnscheduledViewData
		taskDayData          views.TaskGroupViewData
	)
	// Check for an active session
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, r, fmt.Errorf("%s: %w %s", caller, ErrInvalidSessionID, sessionID))
		return
	}

	if user, err = h.services.GetUserBySessionID(sessionID); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}
	if user.ID == 0 {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, ErrEmptySessionID))
		return
	}

	if taskID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	dayString = r.URL.Query().Get("due-date")

	// Set to unscheduled
	if dayString == "" {
		if updatedTask, err = h.services.UpdateDueDate(taskID, user.ID, time.Time{}); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
			return
		}

		if tasks, err = h.services.GetUnscheduledTasks(user.ID); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
			return
		}

		unscheduledData = h.parseTaskListToUnscheduledTasks(tasks)

		if err = views.UnscheduledTasks(unscheduledData).Render(r.Context(), w); err != nil {
			h.HandleError(w, r, fmt.Errorf("%s: %w UnscheduledTasks: %s", caller, ErrRenderingTemplate, err))
			return
		}

		return
	}

	// Grab the day of tasks for the new due date
	if newDueDay, err = h.parseDay(dayString); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if updatedTask, err = h.services.UpdateDueDate(taskID, user.ID, newDueDay); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	if tasks, err = h.services.GetTasksForDay(user.ID, updatedTask.DueDate); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w", caller, err))
		return
	}

	taskDayData = h.parseTaskListToTaskGroupViewData(tasks)

	if err = views.DayTaskList(taskDayData).Render(r.Context(), w); err != nil {
		h.HandleError(w, r, fmt.Errorf("%s: %w DayTaskList: %s", caller, ErrRenderingTemplate, err))
		return
	}
}
