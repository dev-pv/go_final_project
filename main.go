package main

import (
	"log"
	"net/http"
	"os"

	"todo-app/internal/db"
)

func main() {
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	_, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
