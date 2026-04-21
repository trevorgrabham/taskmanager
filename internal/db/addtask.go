package db

import (
	"database/sql"
	"fmt"
	"local/taskmanager/internal/task"
)

func AddTask(db *sql.DB, newTask task.Task) (id int, err error) {
	var (
		queryColumns, queryPlaceholders []string
		queryArgs                       []any
		tx                              *sql.Tx
	)
	if err = checkDBConnection(db); err != nil {
		return 0, err
	}
	if err = checkTaskNotEmpty(newTask); err != nil {
		return 0, err
	}

	if tx, err = startDBSession(db); err != nil {
		return 0, err
	}

	if queryColumns, queryPlaceholders, queryArgs, err = setupQueryParams(tx, newTask); err != nil {
		return 0, err
	}

	fmt.Println(queryColumns)
	fmt.Println(queryArgs...)

	if id, err = addNewTask(tx, queryColumns, queryPlaceholders, queryArgs); err != nil {
		return 0, err
	}

	_ = tx.Commit()
	return id, nil
}
