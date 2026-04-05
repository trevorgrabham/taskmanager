package main

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views"
	"log"
	"net/http"
	"slices"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		db, err := sqlite.Connect()
		if err != nil {
			log.Fatal(err)
		}
		var tasks map[int64]task.TaskList
		tasks, err = sqlite.ListWeeklyTasks(db, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		keys := make([]int64, 0, 7)
		for k := range tasks {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		err = views.Layout(views.Index(tasks, keys)).Render(r.Context(), w)
		if err != nil {
			log.Fatal(err)
		}
	})

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
