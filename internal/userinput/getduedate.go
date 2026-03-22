package userinput

import (
	"bufio"
	"fmt"
	"local/taskmanager/internal/parse"
	"local/taskmanager/internal/task"
	"os"
	"strings"
)

func promptDueDate(prompt string) (task.TaskDueDate, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	reader := bufio.NewReader(os.Stdin)
	dueDateString, err := reader.ReadString('\n')
	if err != nil {
		return task.TaskDueDate{}, err
	}

	dueDateString = strings.TrimSpace(strings.ToLower(dueDateString))
	if dueDateString == "" || dueDateString == "n" || dueDateString == "no" || dueDateString == "none" {
		return task.TaskDueDate{}, nil
	}

	var dueDate task.TaskDueDate
	dueDate, err = parse.ParseDateAndTime(dueDateString)
	if err != nil {
		return task.TaskDueDate{}, err
	}

	return dueDate, nil
}
