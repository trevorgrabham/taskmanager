package main

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/task"
	"local/taskmanager/internal/views"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := views.Layout(views.Index()).Render(r.Context(), w)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/todays-tasks", func(w http.ResponseWriter, r *http.Request) {
		db, err := sqlite.Connect()
		if err != nil {
			log.Fatal(err)
		}
		var tasks task.TaskList
		tasks, err = sqlite.ListDailyTasks(db)
		if err != nil {
			log.Fatal(err)
		}
		err = views.List(tasks).Render(r.Context(), w)
		if err != nil {
			log.Fatal(err)
		}
	})

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
