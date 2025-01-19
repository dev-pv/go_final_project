package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"todo-app/internal/app"
)

func DoneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if err := app.MarkTaskDone(db, id); err != nil {
			log.Printf("MarkTaskDone Error: %v", err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Write([]byte(`{}`))
	}
}
