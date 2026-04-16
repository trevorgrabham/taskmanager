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
		info = views.PageInfo{StartDay: time.Now()}
		err  error
	)
	// Tasks for the upcoming week
	info.WeekOfTasks, err = sqlite.ListWeekOfTasks(h.DB, info.StartDay)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}

	// Get the week of tasks in order
	tasksPerCategory := make(map[string]task.TaskList)
	for dateKey := range info.WeekOfTasks {
		for catKey := range info.WeekOfTasks[dateKey] {
			for _, t := range info.WeekOfTasks[dateKey][catKey] {
				tasksPerCategory[t.Category] = append(tasksPerCategory[t.Category], t)
			}
		}
	}

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
	var overdue task.TaskList
	overdue, err = sqlite.OverdueTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}
	info.OverdueTasks = make(map[string]task.TaskList)
	for _, t := range overdue {
		info.OverdueTasks[t.Category] = append(info.OverdueTasks[t.Category], t)
	}

	var unscheduled task.TaskList
	unscheduled, err = sqlite.UnscheduledTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}
	info.UnscheduledTasks = make(map[string]task.TaskList)
	for _, t := range unscheduled {
		info.UnscheduledTasks[t.Category] = append(info.UnscheduledTasks[t.Category], t)
	}

	w.Header().Set("Content-Type", "text/html")
	err = shared.Layout(dashboard.Index(info)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("index: %s\n", err)
		return
	}
}
