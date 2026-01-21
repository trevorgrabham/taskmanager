package main

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
	"local/taskmanager2.0/internal/userinput"
	"log"
	"os"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, `taskmanager - manages upcoming tasks
Usage:
  taskmanager [-h] [-a] [-c] [-d] [-l] [-s]
  -h print help
  -a add a new task 
  -c complete a task 
  -d remove a task
  -l list incomplete tasks
  -s search for tasks

`)
}

func completeTaskMenu(tasks task.TaskList) error {
	selection, err := userinput.GetTaskSelection("Which task would you like to complete?", tasks.Incomplete())
	if err != nil {
		return err
	}
	if selection == nil {
		return nil
	}

	selection.Done = true
	selection.CompletionDate = task.TaskDueDate(time.Now())
	return nil
}

func main() {
	tasks, err := task.UnmarshalTasks()
	if err != nil {
		log.Fatal(err)
	}

	var didSearch, printedUsage bool
	for _, flag := range os.Args[1:] {
		switch flag {
		case "-h":
			usage()
			printedUsage = true
		case "-a":
			var newTask task.Task
			newTask, err = userinput.AddTaskMenu()
			if err != nil {
				log.Fatal(err)
			}

			tasks = append(tasks, &newTask)
			fmt.Println()
			fmt.Println()
		case "-c":
			err = completeTaskMenu(tasks.Incomplete())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println()
		case "-d":
		case "-p":
			err = userinput.PushTaskMenu(tasks.Incomplete())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println()
		case "-s":
			err = userinput.SearchTaskMenu(tasks)
			if err != nil {
				log.Fatal(err)
			}
			
			didSearch = true
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n\n", flag)
			usage()
			printedUsage = true
		}
	}

	if !didSearch && !printedUsage { fmt.Println(tasks.Incomplete().Sort()) }

	outFile, err := os.OpenFile(task.SaveFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = task.MarshalTasks(outFile, tasks...)
	if err != nil {
		log.Fatal(err)
	}
}
