package views

import (
	"local/taskmanager/internal/task"
	"time"
)

type PageInfo struct {
	StartDay         time.Time
	WeekOfTasks      map[int]map[string]task.TaskList
	OverdueTasks     map[string]task.TaskList
	FavCategoryTasks task.TaskList
	UnscheduledTasks map[string]task.TaskList
}
