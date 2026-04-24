package views

import "time"

type TaskFormViewData struct {
	TaskID              int
	UserID              int
	Title               string
	Description         string
	Category            string
	CategorySuggestions []string
	DueDateDay          string
	DueDateTime         string
	RecurringValue      string
	RecurringUnit       string
}

type DashboardViewData struct {
	StartDay         time.Time
	WeekOfTasks      map[int]TaskGroupViewData
	OverdueTasks     OverdueViewData
	UnscheduledTasks UnscheduledViewData
}

type OverdueViewData struct {
	Categories map[string]CategoryViewData
}

type UnscheduledViewData struct {
	Categories map[string]CategoryViewData
}

type TaskGroupViewData struct {
	Date       time.Time
	Categories map[string]CategoryViewData
}

type CategoryViewData struct {
	Category string
	Tasks    []TaskViewData
}

type TaskViewData struct {
	ID    int
	Title string
	Done  bool
}
