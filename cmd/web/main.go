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

	http.HandleFunc("/", handlers.SessionCookieMiddleware(handlers.IndexHandler))
	http.HandleFunc("/add-task", handlers.SessionCookieMiddleware(handlers.AddTaskHandler))
	http.HandleFunc("/complete", handlers.SessionCookieMiddleware(handlers.CompleteTaskHandler))
	http.HandleFunc("/edit-task", handlers.SessionCookieMiddleware(handlers.EditTaskHandler))
	http.HandleFunc("/delete", handlers.SessionCookieMiddleware(handlers.DeleteTaskHandler))
	http.HandleFunc("/login", handlers.SessionCookieMiddleware(handlers.LoginHandler))
	http.HandleFunc("/signup", handlers.SessionCookieMiddleware(handlers.SignupHandler))
	http.HandleFunc("/task", handlers.SessionCookieMiddleware(handlers.TaskInfoHandler))
	http.HandleFunc("/toggle-complete", handlers.SessionCookieMiddleware(handlers.ToggleCompleteHandler))
	http.HandleFunc("/toggle-recurring", handlers.SessionCookieMiddleware(handlers.ToggleRecurringTaskHandler))
	http.HandleFunc("/toggle-time", handlers.SessionCookieMiddleware(handlers.ToggleTimeHandler))
	http.HandleFunc("/update-task", handlers.SessionCookieMiddleware(handlers.ParseTaskForm))
	http.HandleFunc("/update-task-duedate", handlers.SessionCookieMiddleware(handlers.UpdateTaskDueDateHandler))

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
