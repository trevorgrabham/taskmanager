package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	sqlite "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
)

func SearchTask(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("searching task: cannot search for a task with a nil database")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Search task titles:\t")
	userInput, err := reader.ReadString('\n')
	if err != nil { return fmt.Errorf("searching task: %s", err) }

	userInput = strings.TrimSpace(strings.ToLower(userInput))
	var matches task.TaskList
	matches, err = sqlite.SearchTasks(db, userInput)
	if err != nil { return fmt.Errorf("searching task: %s", err) }

	fmt.Println()
	fmt.Println(matches)
	return nil
}

func SearchTaskMetaData(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("searching task: cannot search for a task with a nil database")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Search task titles:\t")
	userInput, err := reader.ReadString('\n')
	if err != nil { return fmt.Errorf("searching task: %s", err) }

	userInput = strings.TrimSpace(strings.ToLower(userInput))
	var matches []task.TaskMetaData
	matches, err = sqlite.SearchTasksMetaData(db, userInput)
	if err != nil { return fmt.Errorf("searching task: %s", err) }

	fmt.Println()
	fmt.Println(matches)
	return nil
}
