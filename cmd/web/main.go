package main

import (
	sqlite "local/taskmanager/db"
	"local/taskmanager/middleware"
	"local/taskmanager/routes"
	"local/taskmanager/services"
	"log"
	"net/http"
)

func main() {
	repo, err := sqlite.NewRepo() 
	if err != nil { log.Fatal(err) }
	handler := routes.NewHandler(services.NewService(&repo))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", middleware.SessionIDCookieMiddleware(handler.IndexHandler))
	http.HandleFunc("/add-task", middleware.SessionIDCookieMiddleware(handler.AddTaskHandler))
	http.HandleFunc("/complete", middleware.SessionIDCookieMiddleware(handler.CompleteTaskHandler))
	http.HandleFunc("/edit-task", middleware.SessionIDCookieMiddleware(handler.EditTaskHandler))
	http.HandleFunc("/delete", middleware.SessionIDCookieMiddleware(handler.DeleteTaskHandler))
	http.HandleFunc("/login", middleware.SessionIDCookieMiddleware(handler.LoginHandler))
	http.HandleFunc("/signup", middleware.SessionIDCookieMiddleware(handler.SignupHandler))
	http.HandleFunc("/task", middleware.SessionIDCookieMiddleware(handler.TaskInfoHandler))
	http.HandleFunc("/toggle-complete", middleware.SessionIDCookieMiddleware(handler.ToggleCompleteHandler))
	http.HandleFunc("/toggle-recurring", middleware.SessionIDCookieMiddleware(handler.ToggleRecurringTaskHandler))
	http.HandleFunc("/toggle-time", middleware.SessionIDCookieMiddleware(handler.ToggleTimeHandler))
	http.HandleFunc("/update-task", middleware.SessionIDCookieMiddleware(handler.ParseTaskForm))
	http.HandleFunc("/update-task-duedate", middleware.SessionIDCookieMiddleware(handler.UpdateTaskDueDateHandler))

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
