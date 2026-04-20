package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
	"time"
)

func ListDailyTasks(db *sql.DB, day time.Time) (task.TaskList, error) {
	var (
		err error 
		start, end time.Time
	)
	if err = checkDBConnection(db); err != nil { return nil, err }
	if err = checkDateInitialized(day); err != nil { return nil, err }

	start, end = computeStartAndEndDays(day, 1)

	return queryTasks(db, taskQueryParams{
		WhichTasks: incomplete,
		From: start,
		To: end,
	})
}
