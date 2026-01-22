package userinput

import (
	"errors"
	"fmt"
	"local/taskmanager2.0/internal/task"
	"time"
)

func PushTaskMenu(tasks task.TaskList) error {
	fmt.Println("Which task would you like to push?")
	for i := range tasks {
		fmt.Printf("%d. %.80s\n", i+1, tasks[i].Title)
	}
	fmt.Printf("0. Cancel\n")
	var selection, numDays int
	_, err := fmt.Scan(&selection)
	if err != nil {
		return err
	}

	if selection > len(tasks) || selection < 0 {
		return fmt.Errorf("%d is not a valid selection", selection)
	}
	if selection == 0 {
		return nil
	}

	fmt.Print("How many days should it be pushed back?\t")
	_, err = fmt.Scan(&numDays)
	if err != nil {
		return err
	}

	if numDays < 1 {
		return errors.New("cannot push back by a negative number of days")
	}

	taskDueDate := task.TaskDueDate(time.Time(*tasks[selection-1].DueDate).AddDate(0, 0, numDays))
	tasks[selection-1].DueDate = &taskDueDate

	if time.Time(taskDueDate).After(time.Now()) {
		tasks[selection-1].Type = task.Upcoming
	} else {
		tasks[selection-1].Type = task.Due
	}

	return nil
}
