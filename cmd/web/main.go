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

	handlers := routes.Handlers{DB: db}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/add-task-form", handlers.AddTaskHandler)
	http.HandleFunc("/add-task", handlers.ParseAddTaskHandler)
	http.HandleFunc("/recurring-toggle", handlers.ToggleRecurringTaskHandler)

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
