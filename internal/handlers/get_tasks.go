package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"todo-app/internal/app"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		tasks, err := app.GetTasks(db)
		if err != nil {
			log.Printf("GetTasks Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := map[string]interface{}{
			"tasks": tasks,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Encode Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to encode tasks"})
			return
		}
	}
}
