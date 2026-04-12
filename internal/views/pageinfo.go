package views

import (
	"local/taskmanager/internal/task"
	"time"
)

type PageInfo struct {
	StartDay                time.Time
	WeekOfTasks             map[int]map[string]task.TaskList
	WeekCategoryOrderedKeys map[int][]string
	OverdueTasks            task.TaskList
	FavCategoryTasks        task.TaskList
	UnscheduledTasks        task.TaskList
}
