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

type Handlers struct {
	services Services
}

func NewHandler(s Services) Handlers { return Handlers{services: s} }

type Services interface {
	AddTask(t services.Task, userID int) (addedTask services.Task, err error)
	CompleteTask(taskID, userID int) error
	DeleteTask(taskID, userID int) error
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
	UpdateTask(t services.Task, userID int) (updatedTask services.Task, err error)
}

// To parse taskform.templ
type TaskFormData struct {
	SessionID      string
	ID             string
	Title          string
	Description    string
	Category       string
	DueDateDay     string
	DueDateTime    string
	RecurringValue string
	RecurringUnit  string
}

// ============================== Translate to Service Data ==========================================
func (h Handlers) parseDateAndTime(dayString, timeString string) (dateTime time.Time, err error) {
	if dayString == "" {
		return time.Time{}, nil
	}
	if timeString == "" {
		timeString = "00:00"
	}

	if dateTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", dayString, timeString), time.Local); err != nil {
		return time.Time{}, err
	}

	return dateTime, nil
}

func (h Handlers) parseDay(dayString string) (day time.Time, err error) {
	if dayString == "" {
		return time.Time{}, nil
	}

	day, err = time.Parse("2006-01-02", dayString)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing day from %s: %w", dayString, err)
	}

	return day, nil
}

func (h Handlers) parseID(idString string) (id int, err error) {
	if idString == "" {
		return -1, nil
	}

	id, err = strconv.Atoi(idString)
	if err != nil {
		return -1, fmt.Errorf("parsing ID %s: %w", idString, err)
	}

	return id, nil
}

func (h Handlers) parseRecurringPeriod(value, unit string) (recurringPeriod string, err error) {
	if value == "" || unit == "" {
		return "", nil
	}

	if unit != "days" && unit != "weeks" && unit != "months" {
		return "", fmt.Errorf("%s is not a valid recurring unit", unit)
	}

	var valueCheck int
	if valueCheck, err = strconv.Atoi(value); err != nil {
		return "", fmt.Errorf("%s is not a valid recurring value", value)
	}

	return fmt.Sprintf("%d %s", valueCheck, unit), nil
}

// ============================== Translate to ViewData ==========================================

func (h Handlers) parseTaskListToOverdueTasks(tasks services.TaskList) (sortedTasks views.OverdueViewData) {
	sortedTasks.Categories = make(map[string]views.CategoryViewData)
	for _, t := range tasks {
		newTaskViewData := views.TaskViewData{
			ID:    t.ID,
			Title: t.Title,
			Done:  t.Done,
		}
		updatedTasks := append(sortedTasks.Categories[t.Category].Tasks, newTaskViewData)
		sortedTasks.Categories[t.Category] = views.CategoryViewData{Category: t.Category, Tasks: updatedTasks}
	}

	return sortedTasks
}

func (h Handlers) parseTaskToTaskInfoViewData(t services.Task) (taskInfoData taskinfo.TaskInfoViewData) {
	taskInfoData.ID = t.ID
	taskInfoData.Title = t.Title
	taskInfoData.Category = t.Category
	taskInfoData.Description = t.Description
	taskInfoData.DueDate = t.DueDate
	taskInfoData.RecurringPeriod = t.RecurringPeriod

	return taskInfoData
}

func (h Handlers) parseTaskToTaskViewData(taskData services.Task) (taskViewData views.TaskViewData) {
	taskViewData.ID = taskData.ID
	taskViewData.Title = taskData.Title
	taskViewData.Done = taskData.Done

	return taskViewData
}

func (h Handlers) parseTaskToTaskFormViewData(taskData services.Task) (taskFormViewData views.TaskFormViewData, err error) {
	var split []string
	taskFormViewData.TaskID = taskData.ID
	taskFormViewData.UserID = taskData.UserID
	taskFormViewData.Title = taskData.Title
	taskFormViewData.Category = taskData.Category
	taskFormViewData.Description = taskData.Description
	if taskData.RecurringPeriod != "" {
		split = strings.Split(taskData.RecurringPeriod, " ")
		if len(split) != 2 {
			return views.TaskFormViewData{}, fmt.Errorf("unknown recurring period format %s", taskData.RecurringPeriod)
		}
		if _, err = strconv.Atoi(split[0]); err != nil {
			return views.TaskFormViewData{}, fmt.Errorf("unknown recurring period value %v", split[0])
		}
		switch split[1] {
		case "days", "weeks", "months":
		// do nothing
		default:
			return views.TaskFormViewData{}, fmt.Errorf("unknown recurring period unit %v", split[1])
		}

		taskFormViewData.RecurringValue, taskFormViewData.RecurringUnit = split[0], split[1]
	}

	if !taskData.DueDate.IsZero() {
		dateString := time.Time(taskData.DueDate).Format("2006-01-02 15:04")
		split = strings.Split(dateString, " ")
		if len(split) != 2 {
			return views.TaskFormViewData{}, fmt.Errorf("unknown due date format %v", taskData.DueDate)
		}
		taskFormViewData.DueDateDay, taskFormViewData.DueDateTime = split[0], split[1]
	}

	return taskFormViewData, nil
}

func (h Handlers) parseTaskListToTaskGroupViewData(tasks services.TaskList) (sortedTasks views.TaskGroupViewData) {
	sortedTasks.Categories = make(map[string]views.CategoryViewData)
	sortedTasks.Date = tasks[0].DueDate
	for _, t := range tasks {
		newTaskViewData := views.TaskViewData{
			ID:    t.ID,
			Title: t.Title,
			Done:  t.Done,
		}
		updatedTasks := append(sortedTasks.Categories[t.Category].Tasks, newTaskViewData)
		sortedTasks.Categories[t.Category] = views.CategoryViewData{Category: t.Category, Tasks: updatedTasks}
	}

	return sortedTasks
}

func (h Handlers) parseTaskListToUnscheduledTasks(tasks services.TaskList) (sortedTasks views.UnscheduledViewData) {
	sortedTasks.Categories = make(map[string]views.CategoryViewData)
	for _, t := range tasks {
		newTaskViewData := views.TaskViewData{
			ID:    t.ID,
			Title: t.Title,
			Done:  t.Done,
		}
		updatedTasks := append(sortedTasks.Categories[t.Category].Tasks, newTaskViewData)
		sortedTasks.Categories[t.Category] = views.CategoryViewData{Category: t.Category, Tasks: updatedTasks}
	}

	return sortedTasks
}

func (h Handlers) parseTaskListToWeekOfTasks(tasks services.TaskList, startDay time.Time) (sortedTasks map[int]views.TaskGroupViewData) {
	sortedTasks = make(map[int]views.TaskGroupViewData)
	for i := range 7 {
		if i != 0 {
			startDay = startDay.AddDate(0, 0, 1)
		}
		sortedTasks[int(startDay.Weekday())] = views.TaskGroupViewData{Date: startDay, Categories: make(map[string]views.CategoryViewData)}
	}
	fmt.Printf("WeekOfTasks map intialized\n\n%v\n", sortedTasks)
	for _, t := range tasks {
		weekday := int(t.DueDate.Weekday())
		newTaskViewData := views.TaskViewData{ID: t.ID, Title: t.Title, Done: t.Done}

		if sortedTasks[weekday].Date.IsZero() {
			sortedTasks[weekday] = views.TaskGroupViewData{
				Date:       time.Time(t.DueDate),
				Categories: make(map[string]views.CategoryViewData),
			}
		}

		updatedTasks := append(sortedTasks[weekday].Categories[t.Category].Tasks, newTaskViewData)
		sortedTasks[weekday].Categories[t.Category] = views.CategoryViewData{Category: t.Category, Tasks: updatedTasks}
	}

	fmt.Printf("WeekOfTasks map finished\n\n%v\n", sortedTasks)
	return sortedTasks
}
