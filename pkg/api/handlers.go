package api

import (
	"encoding/json"
	"go_final_project-main/pkg/db"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	var tasks []*db.Task

	limit, err := strconv.Atoi(os.Getenv("TODO_TASKS_LIMIT"))
	if err != nil {
		limit = 50
	}

	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if search != "" {
		tasks, err = db.SearchTasks(search, limit)
	} else {
		tasks, err = db.Tasks(limit)
	}

	if err != nil {
		log.New(w, "failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		writeJSONError(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResponse{Tasks: tasks}, http.StatusOK)
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var now time.Time
	var err error
	if r.FormValue("now") == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, r.FormValue("now"))
		if err != nil {
			http.Error(w, "Invalid now parameter", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		return
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := checkTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.New(w, "failed to add task to database: "+err.Error(), http.StatusInternalServerError)
		writeJSONError(w, "Failed to add task to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	writeJSON(w, response, http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.New(w, "failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		writeJSONError(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSONError(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	if err := checkTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		log.New(w, "failed to update task: "+err.Error(), http.StatusInternalServerError)
		writeJSONError(w, "Failed to update task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		log.New(w, "failed to delete task: "+err.Error(), http.StatusInternalServerError)
		writeJSONError(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		log.New(w, "failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		writeJSONError(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	now := time.Now()
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			log.New(w, "failed to delete task: "+err.Error(), http.StatusInternalServerError)
			writeJSONError(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "Failed to calculate next date: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			log.New(w, "failed to update task date: "+err.Error(), http.StatusInternalServerError)
			writeJSONError(w, "Failed to update task date: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}
