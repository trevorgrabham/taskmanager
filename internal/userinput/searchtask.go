package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager2.0/internal/db"
	"local/taskmanager2.0/internal/task"
	"os"
	"strings"
	"time"
)

func SearchTask(db *sql.DB) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("(i)NCOMPLETE, (c)OMPLETE, or (a)LL?\t")
	incOrComp, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("searching tasks: %s", err)
	}

	var taskType task.WhichTasks
	switch strings.TrimSpace(incOrComp) {
	case "i", "I", "INCOMPLETE", "INC", "Incomplete", "Inc", "incomplete", "inc":
		taskType = task.Inc
	case "c", "C", "COMPLETE", "COMP", "COM", "Complete", "Comp", "Com", "complete", "comp", "com":
		taskType = task.Comp
	default:
		taskType = task.All
	}

	fmt.Print("From? DD/MM[/YY] (default: beginning):\t")
	fromString, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("searching tasks: %s", err)
	}

	fromString = strings.TrimSpace(fromString)
	var from time.Time
	if fromString != "" {
		from, err = time.Parse("02/01", fromString)
		if err != nil {
			from, err = time.Parse("02/01/06", fromString)
			if err != nil {
				return fmt.Errorf("searching tasks: %s", err)
			}
		}
		var year int
		if from.Year() == 0 {
			year = time.Now().Year()
		} else {
			year = from.Year()
		}
		from = time.Date(year, from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	}

	fmt.Print("To? DD/MM[/YY] (default: end):\t")
	toString, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("searching tasks: %s", err)
	}

	toString = strings.TrimSpace(toString)
	var to time.Time
	if toString != "" {
		to, err = time.Parse("02/01", toString)
		if err != nil {
			to, err = time.Parse("02/01/06", toString)
			if err != nil {
				return fmt.Errorf("searching tasks: %s", err)
			}
		}
		var year int
		if to.Year() == 0 {
			year = time.Now().Year()
		} else {
			year = to.Year()
		}
		to = time.Date(year, to.Month(), to.Day(), 0, 0, 0, 0, to.Location())
	} else {
		to = time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	fmt.Print("Category? (default: all):\t")
	category, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("searching tasks: %s", err)
	}

	category = strings.TrimSpace(category)
	var tasks task.TaskList
	tasks, err = sqlite.QueryTasks(db, sqlite.QueryParams{WhichTasks: taskType, Category: category, From: from, To: to})
	if err != nil {
		return fmt.Errorf("searching tasks: %s", err)
	}

	fmt.Println()
	fmt.Println(tasks)
	return nil
}
