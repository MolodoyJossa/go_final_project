package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project-main/pkg/db"
)

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]string{
		"error": message,
	}
	writeJSON(w, response, statusCode)
}

func checkTask(task *db.Task) error {
	now := time.Now()

	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	if _, err := time.Parse(DateFormat, task.Date); err != nil {
		return fmt.Errorf("invalid date format")
	}

	if task.Date >= now.Format(DateFormat) {
		return nil
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
	} else {
		task.Date = now.Format(DateFormat)
	}

	return nil
}
