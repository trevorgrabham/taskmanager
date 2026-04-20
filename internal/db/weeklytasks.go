package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
	"time"
)

func ListWeekOfTasks(db *sql.DB, day time.Time) (map[int]map[string]task.TaskList, error) {
	var (
		err        error
		start, end time.Time
		tasks      task.TaskList
	)
	if err = checkDBConnection(db); err != nil {
		return nil, err
	}
	if err = checkDateInitialized(day); err != nil {
		return nil, err
	}

	start, end = computeStartAndEndDays(day, 7)

	if tasks, err = queryTasks(db, taskQueryParams{
		WhichTasks: incomplete,
		From:       start,
		To:         end,
	}); err != nil {
		return nil, err
	}

	return sortWeekOfTasks(tasks), nil
}
