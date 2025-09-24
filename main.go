package main

import (
	"log"

	"github.com/joho/godotenv"

	"go_final_project-main/server"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env не найден или не загружен")
	}

	server.Start()
}
