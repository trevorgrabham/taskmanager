package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func ListWeekOfTasks(db *sql.DB, start time.Time) (map[int64]task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("listing weekly: cannot get tasks for a nil database")
	}
	if start.IsZero() {
		return nil, fmt.Errorf("listing weekly: cannot list tasks without a day")
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end := time.Date(start.Year(), start.Month(), start.Day()+6, 23, 59, 0, 0, start.Location())
	rows, err := db.Query(`
		SELECT id, title, category, description, due_date, completion_date, done 
		FROM task 
		WHERE done = 0 AND due_date >= ? AND due_date <= ?
		ORDER BY due_date ASC, category`,
		start.Unix(),
		end.Unix())
	if err != nil {
		return nil, fmt.Errorf("listing weekly: %s", err)
	}
	defer rows.Close()

	var (
		tasks      task.TaskList
		tasksByDay = make(map[int64]task.TaskList)
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
		if i == 0 {
			tasksByDay[start.Unix()] = task.TaskList{}
		}
		tasksByDay[start.AddDate(0, 0, i).Unix()] = task.TaskList{}
	}
	for _, t := range tasks {
		dueDate := time.Time(t.DueDate)
		date := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, time.Local)
		tasksByDay[date.Unix()] = append(tasksByDay[date.Unix()], t)
	}

	return tasksByDay, nil
}
