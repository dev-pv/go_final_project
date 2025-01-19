package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"todo-app/internal/app"
)

// GetTasksHandler обрабатывает GET-запрос /api/tasks.
// Возвращает JSON:
//
//	{"tasks":[
//	    {"id":"123","date":"20240119","title":"...","comment":"","repeat":""},
//	    ...
//	]}
//
// или при ошибке:
//
//	{"error":"..."}
func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tasks, err := app.GetTasks(db)
		if err != nil {
			log.Printf("GetTasks Error: %v", err)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Оборачиваем список в объект {"tasks": [...]}
		resp := map[string]interface{}{
			"tasks": tasks,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Encode Error: %v", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to encode tasks"})
			return
		}
	}
}
