package main

import (
	"fmt"
	sqlitedb "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
	"local/taskmanager2.0/internal/userinput"
	"log"
	"os"
	"time"
)

func main() {
	db, err := sqlitedb.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tasks task.TaskList
	tasks, err = sqlitedb.QueryAll(db)
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
		addLoop:
			for {
				var newTask task.Task
				newTask, err = userinput.AddTaskMenu()
				if err != nil {
					log.Fatal(err)
				}

				tasks = append(tasks, &newTask)
				err = sqlitedb.AddTask(db, newTask)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Anything else to add? (y/n)\t")
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
					break addLoop
				}
			}

			fmt.Println()
			fmt.Println()
		case "-c":
		completeLoop:
			for {
				var selection *task.Task
				selection, err = userinput.GetTaskSelection("Which task would you like to complete?", tasks.Incomplete())
				if err != nil {
					log.Fatal(err)
				}
				if selection == nil {
					break completeLoop
				}

				taskDueDate := task.TaskDueDate(time.Now())
				selection.Done = true
				selection.CompletionDate = &taskDueDate

				fmt.Print("Anything else to complete? (y/n)\t")
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
					break completeLoop
				}
			}

			fmt.Println()
			fmt.Println()
		case "-d":
			err = userinput.DeleteTask(tasks.Incomplete().Sort())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println()
		case "-p":
			err = userinput.PushTaskMenu(tasks.Incomplete().Sort())
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

	if !didSearch && !printedUsage {
		fmt.Println(tasks.Incomplete().Sort())
	}

	outFile, err := os.OpenFile(task.SaveFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = task.MarshalTasks(outFile, tasks...)
	if err != nil {
		log.Fatal(err)
	}

}

func usage() {
	fmt.Fprintf(os.Stderr, `taskmanager - manages upcoming tasks
Usage:
  taskmanager [-h] [-a] [-c] [-d] [-l] [-s]
  -h print help
  -a add a new task 
  -c complete a task 
  -d remove a task
  -s search for tasks

`)
}
