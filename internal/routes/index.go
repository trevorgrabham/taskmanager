package routes

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views"
	"local/taskmanager/internal/views/dashboard"
	"local/taskmanager/internal/views/shared"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

// Renders the index.html page
func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		log.Printf("index: unknown path %s\n", r.URL.Path)
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("index: no database connection provided to handler")
		return
	}

	var (
		info views.PageInfo
		err  error
	)
	// Tasks for the upcoming week
	info.WeekOfTasks, err = sqlite.ListWeekOfTasks(h.DB, time.Now())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}

	// Get the week of tasks in order
	tasksPerCategory := make(map[string]task.TaskList)
	info.WeekDateOrderedKeys = make([]int64, 0, 7)
	info.WeekCategoryOrderedKeys = make(map[int64][]string)
	for dateKey := range info.WeekOfTasks {
		info.WeekCategoryOrderedKeys[dateKey] = make([]string, 0)
		info.WeekDateOrderedKeys = append(info.WeekDateOrderedKeys, dateKey)
		for catKey := range info.WeekOfTasks[dateKey] {
			info.WeekCategoryOrderedKeys[dateKey] = append(info.WeekCategoryOrderedKeys[dateKey], catKey)
			for _, t := range info.WeekOfTasks[dateKey][catKey] {
				tasksPerCategory[t.Category] = append(tasksPerCategory[t.Category], t)
			}
		}
		slices.Sort(info.WeekCategoryOrderedKeys[dateKey])
	}
	slices.Sort(info.WeekDateOrderedKeys)

	// Get the FavCategory (most tasks upcoming)
	var favCategory string
	for k := range tasksPerCategory {
		if favCategory == "" {
			favCategory = k
			continue
		}

		if len(tasksPerCategory[k]) > len(tasksPerCategory[favCategory]) {
			favCategory = k
		}
	}
	info.FavCategoryTasks = tasksPerCategory[favCategory]
	slices.SortFunc(info.FavCategoryTasks, func(a, b task.Task) int {
		return strings.Compare(a.Title, b.Title)
	})

	// Get overdue tasks
	info.OverdueTasks, err = sqlite.OverdueTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}

	info.UnscheduledTasks, err = sqlite.UnscheduledTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = shared.Layout(dashboard.Index(info)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("index: %s\n", err)
		return
	}
}
