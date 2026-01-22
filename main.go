package main

import (
	"fmt"
	"local/taskmanager2.0/internal/task"
	"local/taskmanager2.0/internal/userinput"
	"log"
	"os"
	"time"
)

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
		addLoop:
			for {
				var newTask task.Task
				newTask, err = userinput.AddTaskMenu()
				if err != nil {
					log.Fatal(err)
				}

				tasks = append(tasks, &newTask)

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
		removeLoop:
			for {
				err = userinput.DeleteTask(tasks)
				if err != nil {
					log.Fatal(err)
				}

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

			fmt.Println()
			fmt.Println()
		case "-p":
		pushLoop:
			for {
				err = userinput.PushTaskMenu(tasks.Incomplete())
				if err != nil {
					log.Fatal(err)
				}

				fmt.Print("Anything else to push? (y/n)\t")
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
					break pushLoop
				}
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
  -l list incomplete tasks
  -s search for tasks

`)
}
