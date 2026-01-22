package userinput

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
	"log"
)

func DeleteTask(tasks task.TaskList) error {
removeLoop:
	for {
		fmt.Println("Which task would you like to remove?")
		for i := range tasks {
			fmt.Printf("%-3s \033[%s%.80s\033[0m\n", fmt.Sprintf("%d.", (i+1)), tasks[i].Type.ANSICode(), tasks[i].Title)
		}
		fmt.Printf("\n0.  Cancel\n")
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

		tasks = tasks.FilterDeleted()

		fmt.Print("Anything else to remove? (y/n)\t")
		var more string
		_, err = fmt.Scan(&more)
		if err != nil {
			log.Fatal(err)
		}

		switch more {
		case "y", "Y", "YES", "Yes", "yes":
			fmt.Println()
			fmt.Println()
		default:
			break removeLoop
		}
	}
	return nil
}
