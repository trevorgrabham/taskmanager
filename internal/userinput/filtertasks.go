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

func parseDate(prompt string) (time.Time, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	dateString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, err
	}

	dateString = strings.TrimSpace(dateString)
	if dateString == "" {
		return time.Time{}, nil
	}

	var date time.Time
	date, err = time.Parse("02/01", dateString)
	if err != nil {
		date, err = time.Parse("02/01/06", dateString)
		if err != nil {
			return time.Time{}, err
		}
	}
	var year int
	if date.Year() == 0 {
		year = time.Now().Year()
	} else {
		year = date.Year()
	}
	date = time.Date(year, date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return date, nil
}

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

	from, err := parseDate("From? DD/MM[/YY] (default: beginning):\t")
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}
	queryParams.From = from

	to, err := parseDate("To? DD/MM[/YY] (default: end):\t")
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}
	queryParams.To = to

	var tasks task.TaskList
	tasks, err = sqlite.QueryTasks(db, queryParams)
	if err != nil {
		return fmt.Errorf("filtering tasks: %s", err)
	}

	fmt.Println()
	fmt.Println(tasks)
	return nil
}
