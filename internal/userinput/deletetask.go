package userinput

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
)

func DeleteTask(tasks task.TaskList) error {
	fmt.Println("Which task would you like to remove?")
	for i := range tasks {
		fmt.Printf("%d. %.80s\n", i+1, (tasks)[i].Title)
	}
	fmt.Printf("0. Cancel\n")
	var selection int
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

	*tasks[selection-1] = task.Task{}
	return nil
}
