package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
	"time"
)

func SearchTasks(db *sql.DB, searchKey string) (task.TaskList, error) {
	if db == nil {
		return nil, fmt.Errorf("searching tasks: cannot search a nil database")
	}
	if searchKey == "" {
		return nil, fmt.Errorf("searching tasks: cannot search without a search term")
	}

	rows, err := db.Query(`SELECT id, title, category, description, due_date, completion_date, done, period FROM task LEFT JOIN recurring ON recurring.id = task.recurring_id WHERE title LIKE ?`, "%"+searchKey+"%")
	if err != nil {
		return nil, fmt.Errorf("searching tasks: %s", err)
	}
	defer rows.Close()

	var (
		matches task.TaskList
		t       task.Task
	)
	for rows.Next() {
		t, err = scanTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("searching tasks: %s", err)
		}
		matches = append(matches, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("searching tasks: %s", err)
	}

	return matches, nil
}

func SearchTasksMetaData(db *sql.DB, searchKey string) (task.TaskMetaDataList, error) {
	if db == nil {
		return nil, fmt.Errorf("searching task meta data: cannot search a nil database")
	}
	if searchKey == "" {
		return nil, fmt.Errorf("searching task meta data: cannot search without a search term")
	}

	rows, err := db.Query(`SELECT title, COUNT(*), MAX(completion_date) FROM task WHERE title LIKE ? GROUP BY title`, "%"+searchKey+"%")
	if err != nil {
		return nil, fmt.Errorf("searching task meta data: %s", err)
	}
	defer rows.Close()

	var (
		matches           task.TaskMetaDataList
		name              string
		count             int
		lastCompletedUnix int64
	)
	for rows.Next() {
		err = rows.Scan(&name, &count, &lastCompletedUnix)
		if err != nil {
			return nil, fmt.Errorf("searching task meta data: %s", err)
		}

		matches = append(matches, task.TaskMetaData{TaskName: name, Count: count, LastCompleted: task.TaskDueDate(time.Unix(lastCompletedUnix, 0))})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("searching task meta data: %s", err)
	}

	return matches, nil
}
