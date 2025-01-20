package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"todo-app/internal/app"
)

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed",
			})
			return
		}

		var task app.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Printf("JSON Decode Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "JSON deserialization error",
			})
			return
		}

		id, err := app.AddTask(&task, db)
		if err != nil {
			log.Printf("AddTask Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}
