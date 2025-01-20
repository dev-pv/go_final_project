package main

import (
	"log"
	"net/http"
	"os"

	"todo-app/internal/db"
	"todo-app/internal/handlers"
)

func main() {
	// Порт по умолчанию
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	// Инициализируем БД
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Статика
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Все хэндлеры
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTaskHandler(database)(w, r)
		case http.MethodPost:
			handlers.AddTaskHandler(database)(w, r)
		case http.MethodPut:
			handlers.EditTaskHandler(database)(w, r)
		case http.MethodDelete:
			handlers.DeleteTaskHandler(database)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(database))
	http.HandleFunc("/api/task/done", handlers.DoneTaskHandler(database))

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
