package views

import "local/taskmanager/internal/task"

type PageInfo struct {
	WeekOfTasks             map[int64]map[string]task.TaskList
	WeekDateOrderedKeys     []int64
	WeekCategoryOrderedKeys map[int64][]string
	OverdueTasks            task.TaskList
	FavCategoryTasks        task.TaskList
	UnscheduledTasks        task.TaskList
}
