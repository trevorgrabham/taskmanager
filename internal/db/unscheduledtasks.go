package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func UnscheduledTasks(db *sql.DB) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("unschuduled tasks: cannot get tasks for a nil database")
	}

	rows, err := db.Query(`
		SELECT task.id, title, category, description, due_date, completion_date, done, period
		FROM task 
		LEFT JOIN recurring ON task.recurring_id = recurring.id
		WHERE due_date IS NULL AND done = 0
		ORDER BY category, title`)
	if err != nil {
		return nil, fmt.Errorf("unscheduled tasks: %s", err)
	}
	defer rows.Close()

	var tasks task.TaskList
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("unscheduled tasks: %s", err)
		}

		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unscheduled tasks: %s", err)
	}

	return tasks, nil
}
