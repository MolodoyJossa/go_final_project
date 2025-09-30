package api

import (
	"go_final_project-main/pkg/cfg"
	"net/http"
)

func Init(config *cfg.Config) {
	http.HandleFunc("/api/signin", Signin(config))
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", Auth(config, TaskHandler))
	http.HandleFunc("/api/tasks", Auth(config, TasksHandler))
	http.HandleFunc("/api/task/done", Auth(config, TaskDoneHandler))
}
