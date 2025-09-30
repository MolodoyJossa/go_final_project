package cfg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBFile    string
	JWTSecret []byte
	Password  string
}

func LoadConfig() *Config {
	cfg := &Config{}
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found or not loaded")
	}

	cfg.DBFile = os.Getenv("TODO_DBFILE")
	if cfg.DBFile == "" {
		cfg.DBFile = "scheduler.db"
	}

	cfg.JWTSecret = []byte(os.Getenv("TODO_SECRET"))

	cfg.Password = os.Getenv("TODO_PASSWORD")

	return cfg
}
