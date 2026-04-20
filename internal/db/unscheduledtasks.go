package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
)

func UnscheduledTasks(db *sql.DB) (task.TaskList, error) {
	if err := checkDBConnection(db); err != nil {
		return nil, err
	}

	return queryTasks(db, taskQueryParams{
		WhichTasks:  incomplete,
		Unscheduled: true,
	})
}
