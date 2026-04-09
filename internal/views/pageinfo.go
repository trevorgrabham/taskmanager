package views

import "local/taskmanager/internal/task"

type PageInfo struct {
	WeekOfTasks      map[int64]task.TaskList
	WeekOrderedKeys  []int64
	OverdueTasks     task.TaskList
	FavCategoryTasks task.TaskList
	UnscheduledTasks task.TaskList
}
