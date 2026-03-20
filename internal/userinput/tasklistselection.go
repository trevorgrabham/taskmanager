package userinput

import (
	"fmt"
	"local/taskmanager/internal/task"
)

func GetTaskSelection(prompt string, tasks task.TaskList) (task.Task, error) {
	fmt.Println(prompt)
	for i := range tasks {
		fmt.Printf("%-3s \033[%s%.80s\033[0m\n", fmt.Sprintf("%d.", (i+1)), tasks[i].ANSICode(), tasks[i].Title)
	}
	fmt.Printf("\n0.  Cancel\n")

	var selection int
	_, err := fmt.Scan(&selection)
	if err != nil {
		return task.Task{}, fmt.Errorf("task selection: %s", err)
	}

	if selection > len(tasks) || selection < 0 {
		return task.Task{}, fmt.Errorf("task selection: %d is not a valid selection", selection)
	}
	if selection == 0 {
		return task.Task{}, nil
	}

	return tasks[selection-1], nil
}
