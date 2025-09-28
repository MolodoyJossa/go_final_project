package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/signin", Signin)
	http.HandleFunc("/api/nextdate", Auth(NextDayHandler))
	http.HandleFunc("/api/task", Auth(TaskHandler))
	http.HandleFunc("/api/tasks", Auth(TasksHandler))
	http.HandleFunc("/api/task/done", Auth(TaskDoneHandler))
}
