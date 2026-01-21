package userinput

import (
	"bufio"
	"fmt"
	"local/taskmanager2.0/internal/task"
	"os"
	"strings"
	"time"
)

func SearchTaskMenu(tasks task.TaskList) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("incomplete or complete? (default: both):\t")
	incOrComp, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	switch strings.TrimSpace(incOrComp) {
	case "INCOMPLETE", "INC", "Incomplete", "Inc", "incomplete", "inc":
		tasks = tasks.Incomplete()
	case "COMPLETE", "COMP", "COM", "Complete", "Comp", "Com", "complete", "comp", "com":
		tasks = tasks.Complete()
	}

	fmt.Print("From? DD/MM[/YY] (default: beginning):\t")
	fromString, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	fromString = strings.TrimSpace(fromString)
	var from time.Time
	if fromString != "" {
		from, err = time.Parse("02/01", fromString)
		if err != nil {
			from, err = time.Parse("02/01/06", fromString)
			if err != nil {
				return err
			}
		}
		var year int
		if from.Year() == 0 {
			year = time.Now().Year()
		} else {
			year = from.Year()
		}
		from = time.Date(year, from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		tasks = tasks.From(from)
	}

	fmt.Print("To? DD/MM[/YY] (default: ending):\t")
	toString, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	toString = strings.TrimSpace(toString)
	var to time.Time
	if toString != "" {
		to, err = time.Parse("02/01", toString)
		if err != nil {
			to, err = time.Parse("02/01/06", toString)
			if err != nil {
				return err
			}
		}
		var year int
		if to.Year() == 0 {
			year = time.Now().Year()
		} else {
			year = to.Year()
		}
		to = time.Date(year, to.Month(), to.Day(), 0, 0, 0, 0, to.Location())
		tasks = tasks.To(to)
	}

	fmt.Print("Category? (default: all):\t")
	category, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	category = strings.TrimSpace(category)
	if category != "" {
		tasks = tasks.ByCategory(category)
	}

	fmt.Println()
	fmt.Println(tasks.Sort())
	return nil
}
