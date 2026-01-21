package userinput

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
)

func GetTaskSelection(prompt string, tasks task.TaskList) (*task.Task, error) {
	fmt.Println(prompt)
	for i := range tasks {
		fmt.Printf("%d. %.80s\n", i+1, tasks[i].Title)
	}
	fmt.Printf("0. Cancel\n")
	var selection int
	_, err := fmt.Scan(&selection)
	if err != nil {
		return nil, err
	}

	if selection > len(tasks) || selection < 0 {
		return nil, fmt.Errorf("%d is not a valid selection", selection)
	}
	if selection == 0 {
		return nil, nil
	}

	return tasks[selection-1], nil
}
