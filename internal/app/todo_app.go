package app

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	initialDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	ruleParts := strings.Split(repeat, " ")
	if len(ruleParts) == 0 {
		return "", errors.New("invalid repeat rule")
	}

	switch ruleParts[0] {
	case "d":
		if len(ruleParts) != 2 {
			return "", errors.New("invalid format for daily rule")
		}
		days, err := strconv.Atoi(ruleParts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid days interval")
		}

		for {
			initialDate = initialDate.AddDate(0, 0, days)
			if initialDate.After(now) {
				return initialDate.Format("20060102"), nil
			}
		}

	case "y":
		for {
			initialDate = initialDate.AddDate(1, 0, 0)
			if initialDate.After(now) {
				return initialDate.Format("20060102"), nil
			}
		}

	default:
		return "", errors.New("unsupported repeat rule")
	}
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(task *Task, db *sql.DB) (string, error) {
	if task.Title == "" {
		return "", errors.New("task title is required")
	}

	now := time.Now()
	currentDay := now.Format("20060102")

	if task.Date == "" || task.Date == "today" {
		task.Date = currentDay
	}

	if _, err := time.Parse("20060102", task.Date); err != nil {
		return "", errors.New("invalid date format, expected YYYYMMDD")
	}

	if task.Date < currentDay {
		if task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return "", fmt.Errorf("invalid repeat rule: %v", err)
			}
			task.Date = nextDate
		} else {
			task.Date = currentDay
		}
	}

	if task.Repeat != "" {
		validPrefixes := []string{"d", "y"}
		isValid := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(task.Repeat, prefix) {
				isValid = true
				break
			}
		}
		if !isValid {
			return "", errors.New("unsupported repeat rule")
		}
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM scheduler WHERE date = ? AND title = ?)`
	if err := db.QueryRow(checkQuery, task.Date, task.Title).Scan(&exists); err != nil {
		return "", fmt.Errorf("failed to check existing task: %v", err)
	}
	if exists {
		return "", errors.New("task with the same date and title already exists")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("failed to save task: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve task ID: %v", err)
	}

	return fmt.Sprintf("%d", id), nil
}

func GetTasks(db *sql.DB) ([]Task, error) {
	tasks := make([]Task, 0)

	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50`)
	if err != nil {
		return nil, fmt.Errorf("failed to select tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			t    Task
			id64 int64
		)
		if err := rows.Scan(&id64, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		t.ID = strconv.FormatInt(id64, 10)
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTaskByID(db *sql.DB, id string) (*Task, error) {
	if id == "" {
		return nil, errors.New("no task id provided")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task id '%s'", id)
	}

	query := `SELECT id, date, title, comment, repeat 
              FROM scheduler 
              WHERE id = ?`
	row := db.QueryRow(query, idInt)

	var t Task
	var id64 int64
	err = row.Scan(&id64, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to query task: %v", err)
	}

	t.ID = strconv.FormatInt(id64, 10)
	return &t, nil
}

func EditTask(task *Task, db *sql.DB) error {
	if task.ID == "" {
		return errors.New("empty task id")
	}
	idInt, err := strconv.Atoi(task.ID)
	if err != nil {
		return fmt.Errorf("invalid task id '%s'", task.ID)
	}

	if task.Title == "" {
		return errors.New("task title is required")
	}

	now := time.Now()
	currentDay := now.Format("20060102")

	if task.Date == "" || task.Date == "today" {
		task.Date = currentDay
	}

	if _, err := time.Parse("20060102", task.Date); err != nil {
		return errors.New("invalid date format, expected YYYYMMDD")
	}

	if task.Date < currentDay {
		if task.Repeat != "" {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("invalid repeat rule: %v", err)
			}
			task.Date = nextDate
		} else {
			task.Date = currentDay
		}
	}

	if task.Repeat != "" {
		validPrefixes := []string{"d", "y"}
		isValid := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(task.Repeat, prefix) {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("unsupported repeat rule")
		}
	}

	query := `UPDATE scheduler
              SET date = ?, title = ?, comment = ?, repeat = ?
              WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idInt)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected error: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func MarkTaskDone(db *sql.DB, id string) error {
	if id == "" {
		return errors.New("no id provided")
	}

	t, err := GetTaskByID(db, id)
	if err != nil {
		return fmt.Errorf("MarkTaskDone: %w", err)
	}

	if t.Repeat == "" {
		return DeleteTask(db, id)
	}

	now := time.Now()
	newDate, err := NextDate(now, t.Date, t.Repeat)
	if err != nil {
		return fmt.Errorf("invalid repeat rule: %w", err)
	}
	if err := updateTaskDate(db, id, newDate); err != nil {
		return fmt.Errorf("MarkTaskDone: %w", err)
	}
	return nil
}

func DeleteTask(db *sql.DB, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid task id '%s'", id)
	}

	res, err := db.Exec(`DELETE FROM scheduler WHERE id=?`, idInt)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected error: %w", err)
	}
	if rows == 0 {
		return errors.New("task not found")
	}
	return nil
}

func updateTaskDate(db *sql.DB, id string, newDate string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid task id '%s'", id)
	}

	res, err := db.Exec(`UPDATE scheduler SET date=? WHERE id=?`, newDate, idInt)
	if err != nil {
		return fmt.Errorf("failed to update date: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected error: %w", err)
	}
	if rows == 0 {
		return errors.New("task not found")
	}
	return nil
}
