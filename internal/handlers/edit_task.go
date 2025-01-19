package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"todo-app/internal/app"
)

func EditTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var task app.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Printf("EditTask: JSON Decode Error: %v", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка десериализации JSON"})
			return
		}

		if err := app.EditTask(&task, db); err != nil {
			log.Printf("EditTask Error: %v", err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Write([]byte(`{}`))
	}
}
