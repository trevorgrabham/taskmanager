package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func ListWeekOfTasks(db *sql.DB, start time.Time) (map[int]map[string]task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("listing weekly: cannot get tasks for a nil database")
	}
	if start.IsZero() {
		return nil, fmt.Errorf("listing weekly: cannot list tasks without a day")
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end := time.Date(start.Year(), start.Month(), start.Day()+6, 23, 59, 0, 0, start.Location())
	rows, err := db.Query(`
		SELECT task.id, title, category, description, due_date, completion_date, done, period
		FROM task 
		LEFT JOIN recurring ON task.recurring_id = recurring.id
		WHERE done = 0 AND due_date >= ? AND due_date <= ?
		ORDER BY due_date ASC, category, title`,
		start.Unix(),
		end.Unix())
	if err != nil {
		return nil, fmt.Errorf("listing weekly: %s", err)
	}
	defer rows.Close()

	var (
		tasks      task.TaskList
		tasksByDay = make(map[int]map[string]task.TaskList)
	)
	for rows.Next() {
		var t task.Task
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("listing weekly: %s", err)
		}

		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("listing weekly: %s", err)
	}

	for i := range 7 {
		tasksByDay[i] = make(map[string]task.TaskList)
	}
	for _, t := range tasks {
		dueDate := time.Time(t.DueDate)
		tasksByDay[int(dueDate.Weekday())][t.Category] = append(tasksByDay[int(dueDate.Weekday())][t.Category], t)
	}

	return tasksByDay, nil
}
