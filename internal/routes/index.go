package routes

import (
	"time"
	"log"
	"net/http"
	"slices"
	"local/taskmanager/internal/views"
	"local/taskmanager/internal/task"
	sqlite "local/taskmanager/internal/db"
)

func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		log.Printf("index: unknown path %s", r.URL.Path)
		return
	}
	if h.DB == nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("index: no database connection provided to handler")
		return
	}

	// Tasks for the upcoming week
	var (
		info views.PageInfo
		err  error
	)
	info.WeekOfTasks, err = sqlite.ListWeekOfTasks(h.DB, time.Now())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}

	// Get the week of tasks in order
	tasksPerCategory := make(map[string]task.TaskList)
	info.WeekOrderedKeys = make([]int64, 0, 7)
	for k := range info.WeekOfTasks {
		info.WeekOrderedKeys = append(info.WeekOrderedKeys, k)
		for _, t := range info.WeekOfTasks[k] {
			tasksPerCategory[t.Category] = append(tasksPerCategory[t.Category], t)
		}
	}
	slices.Sort(info.WeekOrderedKeys)

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

	// Get overdue tasks
	info.OverdueTasks, err = sqlite.OverdueTasks(h.DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("index: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = views.Layout(views.Index(info)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("index: %s\n", err)
		return
	}
}
