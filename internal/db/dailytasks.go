package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func ListDailyTasks(db *sql.DB, day time.Time) (task.TaskList, error) {
	if db == nil { return nil, fmt.Errorf("listing daily: cannot get tasks for a nil database") }
	if day.IsZero() { return nil, fmt.Errorf("listing daily: cannot list tasks without a day") }

	cutoff := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 0, 0, day.Location())
	tasks, err := QueryTasks(db, TaskQueryParams{WhichTasks: IncTasks, To: cutoff})
	if err != nil { return nil, fmt.Errorf("listing daily: %s", err) }

	return tasks, nil
}
