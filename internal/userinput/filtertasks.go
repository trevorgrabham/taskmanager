package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/parse"
	"local/taskmanager/internal/task"
	"os"
	"strings"
	"time"
)

func FilterTasks(db *sql.DB, category string) error {
	queryParams := sqlite.TaskQueryParams{Category: category}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("(i)NCOMPLETE, (c)OMPLETE, or (a)LL?\t")
	incOrComp, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}

	switch strings.TrimSpace(strings.ToLower(incOrComp)) {
	case "i", "in", "inc", "incomplete":
		queryParams.WhichTasks = sqlite.IncTasks
	case "c", "com", "comp", "complete":
		queryParams.WhichTasks = sqlite.CompTasks
	default:
		queryParams.WhichTasks = sqlite.AllTasks
	}

	fmt.Print("From? DD/MM[/YY] (default: beginning):\t")
	dateString, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}

	dateString = strings.TrimSpace(dateString)
	if dateString != "" {
		var from time.Time
		from, err = parse.ParseDate(dateString)
		if err != nil {
			return fmt.Errorf("filtering tasks: %s", err)
		}

		queryParams.From = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		fmt.Println(queryParams.From)
	}

	fmt.Print("To? DD/MM[/YY] (default: end):\t")
	dateString, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}

	dateString = strings.TrimSpace(dateString)

	if dateString != "" {
		var to time.Time
		to, err = parse.ParseDate(dateString)
		if err != nil {
			return fmt.Errorf("filtering tasks: %s", err)
		}

		queryParams.To = to
	}

	var tasks task.TaskList
	tasks, err = sqlite.QueryTasks(db, queryParams)
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}

	fmt.Println()
	fmt.Println(tasks)
	return nil
}
