package main

import (
	sqlite "local/taskmanager/internal/db"
	"local/taskmanager/internal/routes"
	"log"
	"net/http"
)

func main() {
	db, err := sqlite.Connect()
	if err != nil {
		log.Fatal("connecting to database: ", err)
	}
	defer db.Close()

	handlers := routes.Handlers{DB: db}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/add-task", handlers.AddTaskHandler)
	http.HandleFunc("/complete", handlers.CompleteTaskHandler)
	http.HandleFunc("/edit-task", handlers.EditTaskHandler)
	http.HandleFunc("/delete", handlers.DeleteTaskHandler)
	http.HandleFunc("/task", handlers.TaskInfoHandler)
	http.HandleFunc("/toggle-complete", handlers.ToggleCompleteHandler)
	http.HandleFunc("/toggle-recurring", handlers.ToggleRecurringTaskHandler)
	http.HandleFunc("/toggle-time", handlers.ToggleTimeHandler)
	http.HandleFunc("/update-task", handlers.ParseTaskForm)
	http.HandleFunc("/update-task-duedate", handlers.UpdateTaskDueDateHandler)

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
