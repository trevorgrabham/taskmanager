package userinput

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
)

func GetTaskSelection(prompt string, tasks task.TaskList) (int, error) {
	fmt.Println(prompt)
	for i := range tasks {
		fmt.Printf("%-3s \033[%s%.80s\033[0m\n", fmt.Sprintf("%d.", (i+1)), tasks[i].Type.ANSICode(), tasks[i].Title)
	}
	fmt.Printf("\n0.  Cancel\n")
	var selection int
	_, err := fmt.Scan(&selection)
	if err != nil {
		return -1, err
	}

	if selection > len(tasks) || selection < 0 {
		return -1, fmt.Errorf("%d is not a valid selection", selection)
	}
	if selection == 0 {
		return -1, nil
	}

	return tasks[selection-1].ID, nil
}
