package routes

import (
	"local/taskmanager/internal/context"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"local/taskmanager/views/dashboard"
	"net/http"
	"time"
)

func (h Handlers) UpdateTaskDueDateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		caller               = "UpdateTaskDueDateHandler"
		err                  error
		sessionID, dayString string
		newDueDay            time.Time
		user                 services.User
		t, updatedTask       services.Task
		tasks                services.TaskList
		unscheduledData      views.UnscheduledViewData
		taskDayData          views.TaskGroupViewData
	)
	if sessionID, err = context.GetSessionID(r.Context()); err != nil {
		h.HandleError(w, NewSessionIDCookieParseError(caller, err))
		return
	}
	if sessionID == "" {
		h.HandleError(w, NewUnauthenticatedError(caller))
		return
	}

	if t.ID, err = h.parseID(r.URL.Query().Get("id")); err != nil {
		h.HandleError(w, NewIDParseError(caller, err))
		return
	}
	if t.ID == -1 {
		h.HandleError(w, NewNoIDError(caller))
		return
	}

	dayString = r.URL.Query().Get("due-date")

	// Set to unscheduled
	if dayString == "" {
		if updatedTask, err = h.services.UpdateDueDate(t.ID, user.ID, time.Time{}); err != nil {
			// TODO: wrap the returned error
			h.HandleError(w, err)
			return
		}

		if tasks, err = h.services.GetUnscheduledTasks(user.ID); err != nil {
			// TODO: wrap the returned error
			h.HandleError(w, err)
			return
		}

		unscheduledData = h.parseTaskListToUnscheduledTasks(tasks)

		if err = dashboard.UnscheduledTasks(unscheduledData).Render(r.Context(), w); err != nil {
			h.HandleError(w, NewTemplateRenderError(caller, "UnscheduledTasks", err))
			return
		}
		return
	}

	// Grab the day of tasks for the new due date
	if newDueDay, err = h.parseDay(dayString); err != nil {
		h.HandleError(w, NewDayParseError(caller, err))
		return
	}

	if updatedTask, err = h.services.UpdateDueDate(t.ID, user.ID, newDueDay); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	if tasks, err = h.services.GetTasksForDay(user.ID, updatedTask.DueDate); err != nil {
		// TODO: wrap the returned error
		h.HandleError(w, err)
		return
	}

	taskDayData = h.parseTaskListToTaskGroupViewData(tasks)

	if err = dashboard.DayTaskList(taskDayData).Render(r.Context(), w); err != nil {
		h.HandleError(w, NewTemplateRenderError(caller, "DayTaskList", err))
		return
	}
}
