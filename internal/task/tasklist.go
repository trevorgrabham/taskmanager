package task

import (
	"slices"
	"sort"
	"strings"
	"time"
)

type TaskList []*Task

func (tl TaskList) Sort() TaskList {
	sort.Slice(tl, func(i, j int) bool {
		if tl[i].Type == tl[j].Type {
			if tl[i].Category == tl[j].Category {
				if tl[i].DueDate == tl[j].DueDate {
					return tl[i].Title < tl[j].Title
				}
				return time.Time(*tl[i].DueDate).Before(time.Time(*tl[j].DueDate))
			}
			return tl[i].Category < tl[j].Category
		}
		return tl[i].Type < tl[j].Type
	})
	return tl
}

func (tl TaskList) From(cutoff time.Time) TaskList {
	filteredTasks := make(TaskList, 0, len(tl))
	for i := range tl {
		if cutoff.After(time.Time(*tl[i].DueDate)) {
			continue
		}

		filteredTasks = append(filteredTasks, tl[i])
	}
	return filteredTasks
}

func (tl TaskList) To(cutoff time.Time) TaskList {
	filteredTasks := make(TaskList, 0, len(tl))
	for i := range tl {
		if cutoff.Before(time.Time(*tl[i].DueDate)) {
			continue
		}

		filteredTasks = append(filteredTasks, tl[i])
	}
	return filteredTasks
}

func (tl TaskList) Incomplete() TaskList {
	filteredTasks := make(TaskList, 0, len(tl))
	for i := range tl {
		if tl[i].Done {
			continue
		}

		filteredTasks = append(filteredTasks, tl[i])
	}
	return filteredTasks
}

func (tl TaskList) Complete() TaskList {
	filteredTasks := make(TaskList, 0, len(tl))
	for i := range tl {
		if !tl[i].Done {
			continue
		}

		filteredTasks = append(filteredTasks, tl[i])
	}
	return filteredTasks
}

func (tl TaskList) Due() TaskList {
	return tl.Incomplete().To(time.Now())
}

func (tl TaskList) Upcoming() TaskList {
	return tl.Incomplete().From(time.Now().Add(time.Minute * 1))
}

func (tl TaskList) ByCategory(category string) TaskList {
	filteredTasks := make(TaskList, 0, len(tl))
	for i := range tl {
		if tl[i].Category != category {
			continue
		}

		filteredTasks = append(filteredTasks, tl[i])
	}
	return filteredTasks
}

func (tl TaskList) String() string {
	taskStrings := make([]string, len(tl))
	for i := range tl {
		taskStrings[i] = tl[i].String()
	}
	return strings.Join(taskStrings, "\n\n")
}

func (tl TaskList) FilterDeleted() TaskList {
	return slices.DeleteFunc(tl, func(t *Task) bool {
		return *t == Task{}
	})
}
