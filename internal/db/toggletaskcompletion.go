package db

import (
	"database/sql"
	"local/taskmanager/internal/task"
)

func ToggleTaskComplete(db *sql.DB, taskToToggle task.Task) (updatedTask, newTask task.Task, err error) {
	var (
		tx                    *sql.Tx
		nextDueDate           task.TaskDueDate
		columns, placeholders []string
		args                  []any
	)
	if err = checkDBConnection(db); err != nil {
		return updatedTask, newTask, err
	}
	if err = checkTaskNotEmpty(taskToToggle); err != nil {
		return updatedTask, newTask, err
	}

	if tx, err = startDBSession(db); err != nil {
		return updatedTask, newTask, err
	}
	defer func() { _ = tx.Rollback() }()

	if updatedTask, err = toggleComplete(tx, taskToToggle.ID); err != nil {
		return task.Task{}, task.Task{}, err
	}

	if updatedTask.RecurringPeriod == "" {
		_ = tx.Commit()
		return updatedTask, task.Task{}, nil
	}

	newTask = updatedTask
	newTask.ID = 0

	// Just completed a recurring task. Add a new one to task table
	if updatedTask.Done {
		if nextDueDate, err = computeNextDueDate(taskToToggle.RecurringPeriod); err != nil {
			return task.Task{}, task.Task{}, err
		}

		newTask.DueDate = nextDueDate

		if columns, placeholders, args, err = setupQueryParams(tx, newTask); err != nil {
			return task.Task{}, task.Task{}, err
		}

		if newTask.ID, err = addNewTask(tx, columns, placeholders, args); err != nil {
			return task.Task{}, task.Task{}, err
		}

		_ = tx.Commit()
		return updatedTask, newTask, nil
	}

	// Just uncompleted a recurring task. Need to update the original, and then remove any tasks that match the title, category, description, recurring_id
	if err = deleteFutureRecurringTask(tx, newTask); err != nil {
		return task.Task{}, task.Task{}, err
	}

	_ = tx.Commit()
	return updatedTask, task.Task{}, nil
}
