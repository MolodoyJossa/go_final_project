package main

import (
	"go_final_project-main/pkg/api"
	"go_final_project-main/pkg/db"
	"go_final_project-main/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	api.Init()

	if err := godotenv.Load(); err != nil {
		log.Println(".env not found or not loaded")
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("DATABASE initialization error: %v", err)
	}

	server.Start()
}
