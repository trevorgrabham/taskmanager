package main

import (
	"fmt"
	sqlite "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
	"local/taskmanager2.0/internal/userinput"
	"log"
	"os"
)

func main() {
	db, err := sqlite.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

				err = sqlite.AddTask(db, newTask)
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
			err = userinput.CompleteTask(db)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println()
		case "-d":
			err = userinput.DeleteTask(db)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println()
		case "-p":
			err = userinput.PushTask(db)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println()
		// case "-s":
		// err = userinput.SearchTaskMenu(tasks)
		// if err != nil {
		// log.Fatal(err)
		// }

		// didSearch = true
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n\n", flag)
			usage()
			printedUsage = true
		}
	}

	if !didSearch && !printedUsage {
		tasks, err := sqlite.QueryIncompleteTasks(db)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tasks)
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
