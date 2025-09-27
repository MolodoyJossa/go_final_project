package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", TaskDoneHandler)
}
