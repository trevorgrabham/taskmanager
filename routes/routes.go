// Package routes provides HTTP request handlers and their supporting types.
package routes

import (
	"fmt"
	"local/taskmanager/services"
	"local/taskmanager/views"
	"local/taskmanager/views/taskinfo"
	"strconv"
	"strings"
	"time"
)

// Handler is a wrapper around a service layer interface. It exposes methods to respond to http requests.
type Handler struct {
	services Services
}

// NewHandler initializes a Handler using s to perform our service logic.
func NewHandler(s Services) Handler { return Handler{services: s} }

// Services defines the service layer functions needed by Handler to respond to requests.
type Services interface {
	AddTask(t services.Task) (addedTask services.Task, err error)
	CompleteTask(taskID, userID int) error
	DeleteTask(taskID, userID int) error
	DeleteSession(sessionID string) error
	GetDashboardData(userID int) (services.Dashboard, error) // Aggregate WeekOfTasks, Unscheduled, Overdue
	GetOverdueTasks(userID int) (services.TaskList, error)   // no need to sort the data
	GetTaskByID(taskID, userID int) (services.Task, error)
	GetTasksForDay(userID int, day time.Time) (services.TaskList, error)
	GetUnscheduledTasks(userID int) (services.TaskList, error) // no need to sort the data
	GetUserBySessionID(sessionID string) (services.User, error)
	LoginUser(username, password string) (sessionID string, err error) // username and password are not "". Does not need to SetSessionIDCookie(), but should start session
	GetUserCategorySuggestions(userID int) ([]string, error)
	GetWeekOfTasks(userID int, day time.Time) (services.TaskList, error) // no need to sort the data
	SignupUser(username, password string) (sessionID string, err error)  // username and password are not "". Does not need to SetSessionIDCookie(), but should start session
	TaskToggleComplete(taskID, userID int) (updatedTask services.Task, err error)
	UpdateDueDate(taskID, userID int, dueDate time.Time) (updatedTask services.Task, err error) // if dueDate.IsZero() remove due date. Otherwise, keep the original time, but replace the day
	UpdateTask(t services.Task) (updatedTask services.Task, err error)
}

// TaskFormData represents the raw data submitted from a TaskForm.
type TaskFormData struct {
	SessionID      string
	ID             string // expected format: "n", n > 1
	Title          string
	Description    string
	Category       string
	DueDateDay     string // expected format: "2006-01-02"
	DueDateTime    string // expected format: "15:04"
	RecurringValue string // expected format: "n", n > 1
	RecurringUnit  string // expected format: "days" | "weeks" | "months"
}

// ============================== Translate to Service Data ==========================================

// parseDateAndTime combines a date string and time string into a time.Time value.
//
// If date string is empty, returns time.Time{}.
// If time string is empty, populates it to "00:00".
// Returns an ErrParsingDate if the combined input is not a valid date-time string.
func (h Handler) parseDateAndTime(dayString, timeString string) (dateTime time.Time, err error) {
	if dayString == "" {
		return time.Time{}, nil
	}
	if timeString == "" {
		timeString = "00:00"
	}

	if dateTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", dayString, timeString), time.Local); err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrParsingDate, fmt.Sprintf("%s %s", dayString, timeString))
	}

	return dateTime, nil
}

// parseDay parses a date string into a time.Time value.
//
// If date string is empty, returns time.Time{}.
// Returns an ErrParsingDay if date string is not a valid date-time string.
func (h Handler) parseDay(dayString string) (day time.Time, err error) {
	if dayString == "" {
		return time.Time{}, nil
	}

	day, err = time.Parse("2006-01-02", dayString)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrParsingDay, dayString)
	}

	return day, nil
}

// parseID parses a string value into an integer.
//
// If id string is empty, returns 0.
// Returns an error if the string cannot be parsed into a valid integer.
func (h Handler) parseID(idString string) (id int, err error) {
	if idString == "" {
		return 0, nil
	}

	id, err = strconv.Atoi(idString)
	if err != nil {
		return -1, fmt.Errorf("%w %s", ErrParsingInt, id)
	}

	return id, nil
}

// parseRecurringPeriod combines a value and unit string of the form "<value> <unit>" (e.g. "2 weeks").
//
// If any of the inputs are empty, an empty string is returned.
// If value does not hold a valid integer, an ErrBadRecurringValue is returned.
// The unit string is not validated, validation is left up to the service layer.
func (h Handler) parseRecurringPeriod(value, unit string) (recurringPeriod string, err error) {
	if value == "" || unit == "" {
		return "", nil
	}

	// TODO: remove this code block once the validation in the service layer is set up
	// if unit != "days" && unit != "weeks" && unit != "months" {
	// return "", fmt.Errorf("%s is not a valid recurring unit", unit)
	// }

	var valueCheck int
	if valueCheck, err = strconv.Atoi(value); err != nil {
		return "", fmt.Errorf("%s is not a valid recurring value", value)
	}

	return fmt.Sprintf("%d %s", valueCheck, unit), nil
}

// ============================== Translate to ViewData ==========================================

// parseTaskListToOverdueTasks maps tasks to OverdueViewData, grouped by category.
func (h Handler) parseTaskListToOverdueTasks(tasks services.TaskList) (sortedTasks views.OverdueViewData) {
	sortedTasks.Categories = make(map[string]*views.CategoryViewData)
	for _, t := range tasks {
		category := sortedTasks.Categories[t.Category]
		if category == nil {
			// Must initialize because we are storing pointers
			category = &views.CategoryViewData{Category: t.Category, Tasks: nil}
		}

		category.Tasks = append(category.Tasks, h.parseTaskToTaskViewData(t))
		sortedTasks.Categories[t.Category] = category
	}

	return sortedTasks
}

// parseTaskToTaskInfoViewData maps a task to TaskInfoViewData.
func (h Handler) parseTaskToTaskInfoViewData(t services.Task) (taskInfoData taskinfo.TaskInfoViewData) {
	taskInfoData.ID = t.ID
	taskInfoData.Title = t.Title
	taskInfoData.Category = t.Category
	taskInfoData.Description = t.Description
	taskInfoData.DueDate = t.DueDate
	taskInfoData.RecurringPeriod = t.RecurringPeriod

	return taskInfoData
}

// parseTaskToTaskViewData maps a task to TaskViewData.
func (h Handler) parseTaskToTaskViewData(taskData services.Task) (taskViewData views.TaskViewData) {
	taskViewData.ID = taskData.ID
	taskViewData.Title = taskData.Title
	taskViewData.Done = taskData.Done

	return taskViewData
}

// parseTaskToTaskFormViewData maps a task to TaskFormViewData.
func (h Handler) parseTaskToTaskFormViewData(taskData services.Task) (taskFormViewData views.TaskFormViewData, err error) {
	var split []string
	taskFormViewData.TaskID = taskData.ID
	taskFormViewData.UserID = taskData.UserID
	taskFormViewData.Title = taskData.Title
	taskFormViewData.Category = taskData.Category
	taskFormViewData.Description = taskData.Description
	if taskData.RecurringPeriod != "" {
		split = strings.Split(taskData.RecurringPeriod, " ")
		taskFormViewData.RecurringValue = split[0]
		taskFormViewData.RecurringUnit = split[1]
		// TODO: remove this code block once the recurring period validation is setup at the service layer. Needs to validate both incoming and outgoing.
		// if len(split) != 2 {
		// return views.TaskFormViewData{}, fmt.Errorf("unknown recurring period format %s", taskData.RecurringPeriod)
		// }
		// if _, err = strconv.Atoi(split[0]); err != nil {
		// return views.TaskFormViewData{}, fmt.Errorf("unknown recurring period value %v", split[0])
		// }
		// switch split[1] {
		// case "days", "weeks", "months":
		// do nothing
		// default:
		// return views.TaskFormViewData{}, fmt.Errorf("unknown recurring period unit %v", split[1])
		// }

		// taskFormViewData.RecurringValue, taskFormViewData.RecurringUnit = split[0], split[1]
	}

	if !taskData.DueDate.IsZero() {
		dateString := taskData.DueDate.Format("2006-01-02 15:04")
		split = strings.Split(dateString, " ")
		taskFormViewData.DueDateDay, taskFormViewData.DueDateTime = split[0], split[1]
	}

	return taskFormViewData, nil
}

// parseTaskListToTaskGroupViewData maps tasks to TaskGroupViewData, grouped by category. It is assumed that all tasks share the same due date.
func (h Handler) parseTaskListToTaskGroupViewData(tasks services.TaskList) (sortedTasks views.TaskGroupViewData) {
	if len(tasks) < 1 {
		return
	}
	sortedTasks.Categories = make(map[string]*views.CategoryViewData)
	sortedTasks.Date = tasks[0].DueDate
	for _, t := range tasks {
		category := sortedTasks.Categories[t.Category]
		if category == nil {
			category = &views.CategoryViewData{Category: t.Category, Tasks: nil}
		}
		category.Tasks = append(category.Tasks, h.parseTaskToTaskViewData(t))
		sortedTasks.Categories[t.Category] = category
	}

	return sortedTasks
}

// parseTaskListToUnscheduledTasks maps tasks to UnscheduledViewData, grouped by category.
func (h Handler) parseTaskListToUnscheduledTasks(tasks services.TaskList) (sortedTasks views.UnscheduledViewData) {
	sortedTasks.Categories = make(map[string]*views.CategoryViewData)
	for _, t := range tasks {
		category := sortedTasks.Categories[t.Category]
		if category == nil {
			category = &views.CategoryViewData{Category: t.Category, Tasks: nil}
		}
		category.Tasks = append(category.Tasks, h.parseTaskToTaskViewData(t))
		sortedTasks.Categories[t.Category] = category
	}

	return sortedTasks
}

// parseTaskListToWeekOfTasks maps tasks to a map[int]*TaskGroupViewData object, grouped by due date, category. tasks are assumed to have due dates that are >= startDay and <= (startDay + 7 days).
func (h Handler) parseTaskListToWeekOfTasks(tasks services.TaskList, startDay time.Time) (sortedTasks map[int]*views.TaskGroupViewData) {
	// initialize sortedTasks
	sortedTasks = make(map[int]*views.TaskGroupViewData)
	sortedTasks[int(startDay.Weekday())] = &views.TaskGroupViewData{Date: startDay, Categories: make(map[string]*views.CategoryViewData)}
	for range 6 {
		startDay = startDay.AddDate(0, 0, 1)
		sortedTasks[int(startDay.Weekday())] = &views.TaskGroupViewData{Date: startDay, Categories: make(map[string]*views.CategoryViewData)}
	}

	for _, t := range tasks {
		weekday := int(t.DueDate.Weekday())

		category := sortedTasks[weekday].Categories[t.Category]
		if category == nil {
			category = &views.CategoryViewData{Category: t.Category, Tasks: nil}
		}
		category.Tasks = append(category.Tasks, h.parseTaskToTaskViewData(t))
		sortedTasks[weekday].Categories[t.Category] = category
	}

	return sortedTasks
}
