package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = "7540"
const webDir = "./web"

func getPort() string {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		return envPort
	}
	return port
}

func Start() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := getPort()

	fmt.Printf("Сервер запущен на http://localhost:%s ...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
