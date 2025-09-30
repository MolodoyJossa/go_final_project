package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var Port int
var DBFile string
var FullNextDate = true
var Search = true
var Token string

func tryLoadDotenv() {
	// Пробуем загрузить .env из разных мест
	paths := []string{".env", "../.env", filepath.Join("..", ".env")}
	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			fmt.Println("[settings] .env загружен из:", p)
			return
		}
	}
	fmt.Println("[settings] .env не найден в стандартных местах")
}

func init() {
	tryLoadDotenv()

	portStr := os.Getenv("TODO_PORT")
	Port = 7540
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil && p > 0 && p < 65536 {
			Port = p
		}
	}
	fmt.Println("[settings] Port:", Port)

	DBFile = os.Getenv("TODO_DBFILE")
	if DBFile == "" {
		DBFile = "scheduler.db"
	}
	fmt.Println("[settings] DBFile:", DBFile)

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		fmt.Println("[settings] TODO_PASSWORD не задан, токен не будет получен")
		return
	}
	url := "http://localhost:" + strconv.Itoa(Port) + "/api/signin"
	body, err := json.Marshal(map[string]string{"password": password})
	if err != nil {
		fmt.Println("[settings] Ошибка сериализации пароля:", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[settings] Ошибка запроса токена:", err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var res struct{ Token string }
	if err := json.NewDecoder(resp.Body).Decode(&res); err == nil && res.Token != "" {
		Token = res.Token
		fmt.Println("[settings] Токен получен")
	} else {
		fmt.Println("[settings] Не удалось получить токен")
	}
}
