package db

import (
	"database/sql"
	"fmt"
)

func UpdateTasksList(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("updating tasks_list: cannot update for a nil database")
	}

	_, err := db.Exec(`UPDATE tasks_list SET times_completed = (SELECT COUNT(*) FROM tasks WHERE tasks.task_id = tasks_list.id AND tasks.done = 1), last_time_completed = (SELECT MAX(completion_date) FROM tasks WHERE tasks.task_id = tasks_list.id)`)
	if err != nil {
		return fmt.Errorf("updating tasks_list: %s", err)
	}

	return nil
}
