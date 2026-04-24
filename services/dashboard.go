package services

type Dashboard struct {
	WeekOfTasks      TaskList
	OverdueTasks     TaskList
	UnscheduledTasks TaskList
}
