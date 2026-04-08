package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func OverdueTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("getting overdue tasks: cannot get tasks for a nil database")
	}

	now := time.Now()
	yesterdayAtMidnight := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 0, 0, time.Local)
	rows, err := db.Query(`
		SELECT task.id, title, category, description, due_date, completion_date, done, period
		FROM task 
		LEFT JOIN recurring ON task.id = recurring.task_id
		WHERE done = 0 AND due_date <= ? AND due_date IS NOT NULL
		ORDER BY category, due_date ASC`,
		yesterdayAtMidnight.Unix())
	if err != nil {
		return nil, fmt.Errorf("getting overdue tasks: %s", err)
	}
	defer rows.Close()

	var tasks task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("getting overdue tasks: %s", err)
		}

		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("getting overdue tasks: %s", err)
	}

	return tasks, nil
}
