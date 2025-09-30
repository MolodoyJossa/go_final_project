package main

import (
	"go_final_project-main/pkg/api"
	"go_final_project-main/pkg/cfg"
	"go_final_project-main/pkg/db"
	"go_final_project-main/pkg/server"
)

func main() {

	config := cfg.LoadConfig()

	api.Init(config)

	db.Init(config.DBFile)
	defer db.Close()

	server.Start()
}
