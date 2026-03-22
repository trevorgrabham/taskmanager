package userinput

import (
	"bufio"
	"database/sql"
	"fmt"
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parsePeriod(userInput string) (string, error) {
	if userInput == "" {
		return "", fmt.Errorf("cannot parse an empty period")
	}

	userInput = strings.TrimSpace(strings.ToLower(userInput))

	re := regexp.MustCompile(`^(\d+)\s*([A-Za-z]+)$`)
	matches := re.FindStringSubmatch(userInput)

	if matches == nil {
		return "", fmt.Errorf("bad period format")
	}

	unit := matches[2]
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}

	var period string
	switch unit {
	case "d", "day", "days":
		period = fmt.Sprintf("%d days", value)
	case "w", "week", "weeks":
		period = fmt.Sprintf("%d weeks", value)
	case "m", "month", "months":
		period = fmt.Sprintf("%d months", value)
	default:
		return "", fmt.Errorf("unrecognized format")
	}

	return period, nil
}

func AddTask(db *sql.DB, category string) error {
	if db == nil {
		return fmt.Errorf("adding task: cannot add a task into a nil database")
	}

	reader := bufio.NewReader(os.Stdin)
addLoop:
	for {
		fmt.Printf("title:\t")
		title, err := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}
		if title == "" {
			return fmt.Errorf("adding task: task needs a title")
		}

		if category == "" {
			fmt.Printf("category:\t")
			category, err = reader.ReadString('\n')
			category = strings.TrimSpace(category)
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		}

		fmt.Printf("description:\t")
		description, err := reader.ReadString('\n')
		description = strings.TrimSpace(description)
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		var dueDate task.TaskDueDate
		dueDate, err = promptDueDate("due date (day of the week (Mon/Monday), or form of DD/MM [XX:XX XM]):\t")
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		var id int
		id, err = sqlite.AddTask(db, task.Task{Title: title, Category: category, Description: description, DueDate: dueDate})
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		var recurring, period string
		fmt.Print("Recurring?\t")
		recurring, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		recurring = strings.TrimSpace(strings.ToLower(recurring))
		switch recurring {
		case "y", "yes":
			fmt.Print("How often?\t")
			recurring, err = reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}

			period, err = parsePeriod(recurring)
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}

			err = sqlite.AddRecurringTask(db, task.RecurringTask{TaskID: id, Period: period})
			if err != nil {
				return fmt.Errorf("adding task: %s", err)
			}
		}

		fmt.Print("Anything else to add? (y/n)\t")
		var more string
		more, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("adding task: %s", err)
		}

		more = strings.TrimSpace(strings.ToLower(more))
		switch more {
		case "n", "no":
			break addLoop
		default:
			fmt.Println()
			fmt.Println()
		}
	}
	return nil
}
