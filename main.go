package main

import (
	"flag"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/userinput"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	db, err := sqlite.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = sqlite.Setup(db)
	if err != nil {
		log.Fatal(err)
	}

	var addTaskFlag, completeTaskFlag, deleteTaskFlag, pushTaskFlag, filterTaskFlag, searchTaskFlag, searchTaskMetaDataFlag bool
	flag.BoolVar(&addTaskFlag, "a", false, "add a task")
	flag.BoolVar(&completeTaskFlag, "c", false, "complete a task")
	flag.BoolVar(&deleteTaskFlag, "d", false, "delete a task")
	flag.BoolVar(&pushTaskFlag, "p", false, "push a task")
	flag.BoolVar(&filterTaskFlag, "f", false, "filter tasks")
	flag.BoolVar(&searchTaskFlag, "s", false, "search for matching tasks")
	flag.BoolVar(&searchTaskMetaDataFlag, "S", false, "search for matching task meta data")
	flag.Parse()
	args := flag.Args()

	if len(args) > 1 {
		log.Fatal("too many args")
	}

	var (
		category    string
		shouldPrint = true
	)
	if len(args) == 1 {
		category = args[0]
	} else {
		var workingDir string
		workingDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		var data []byte
		data, err = os.ReadFile(filepath.Join(workingDir, ".tmconfig"))
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}
		category = strings.TrimSpace(strings.ToLower(string(data)))
	}
	switch {
	case addTaskFlag:
		err = userinput.AddTask(db, category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println()
	case completeTaskFlag:
		err = userinput.CompleteTask(db, category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println()
	case deleteTaskFlag:
		err = userinput.DeleteTask(db, category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println()
	case pushTaskFlag:
		err = userinput.PushTask(db, category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println()
	case filterTaskFlag:
		err = userinput.FilterTasks(db, category)
		if err != nil {
			log.Fatal(err)
		}
		shouldPrint = false
	case searchTaskFlag:
		err = userinput.SearchTask(db)
		if err != nil {
			log.Fatal(err)
		}
		shouldPrint = false
	case searchTaskMetaDataFlag:
		err = userinput.SearchTaskMetaData(db)
		if err != nil {
			log.Fatal(err)
		}
		shouldPrint = false
	default:
		// if no flags and no args, then just print todays tasks
		if len(args) < 1 {
			var tasks task.TaskList
			tasks, err = sqlite.ListDailyTasks(db)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(tasks)
			fmt.Println()
			shouldPrint = false
		}
	}
	if shouldPrint {
		var tasks task.TaskList
		tasks, err = sqlite.QueryTasks(db, sqlite.TaskQueryParams{WhichTasks: sqlite.IncTasks, Category: category})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(tasks)
		fmt.Println()
	}
}

// func usage() {
// fmt.Fprintf(os.Stderr, `taskmanager - manages upcoming tasks
// Usage:
// taskmanager [-h] [-a] [-c] [-d] [-f] [-p] [-s] [-S]
// -h print help
// -a add a new task
// -c complete a task
// -d remove a task
// -f filter tasks
// -p push a task due date
// -s search for tasks
// -S search for task meta data
//
// `)
// }
