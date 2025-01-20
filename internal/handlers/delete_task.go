package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"todo-app/internal/app"
)

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		id := r.URL.Query().Get("id")
		if err := app.DeleteTask(db, id); err != nil {
			log.Printf("DeleteTask error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Write([]byte(`{}`))
	}
}
