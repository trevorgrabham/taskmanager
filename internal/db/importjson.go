package db

import (
	"database/sql"
	"errors"
	"fmt"

	"local/taskmanager2.0/internal/task"
)

func ImportFromJSON(db *sql.DB, fileName string) error {
	if db == nil {
		return errors.New("importing from json: cannot import to a nil db")
	}
	if fileName == "" {
		return errors.New("importing from json: cannot import without a file name")
	}

	tasks, err := task.UnmarshalTasks(fileName)
	if err != nil {
		return fmt.Errorf("unmarshaling json: %s", err)
	}
	if len(tasks) <= 0 {
		return fmt.Errorf("cannot import from %s, has no json data to import", fileName)
	}

	var insertStatement *sql.Stmt
	insertStatement, err = db.Prepare(fmt.Sprintf(`INSERT INTO %s(title, category, description, due_date, completion_date, done) VALUES (?, ?, ?, ?, ?, ?)`, taskTableName))
	if err != nil {
		return fmt.Errorf("preparing insert stmnt: %s", err)
	}
	defer insertStatement.Close()

	for _, t := range tasks {
		err = AddTask(db, *t)
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}
	}
	return nil
}
