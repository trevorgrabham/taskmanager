package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func ListDailyTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil { return nil, fmt.Errorf("listing daily: cannot get tasks for a nil database") }

	now := time.Now()
	cutoff := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, now.Location())
	tasks, err := QueryTasks(db, TaskQueryParams{WhichTasks: IncTasks, To: cutoff})
	if err != nil { return nil, fmt.Errorf("listing daily: %s", err) }

	return tasks, nil
}
