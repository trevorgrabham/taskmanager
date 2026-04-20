package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
	"time"
)

func OverdueTasks(db *sql.DB) (task.TaskList, error) {
	now := time.Now()
	if err := checkDBConnection(db); err != nil {
		return nil, err
	}

	return queryTasks(db, taskQueryParams{
		WhichTasks: incomplete,
		To:         time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 0, 0, time.Local),
	})
}
