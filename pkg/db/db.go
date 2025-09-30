package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT '',
    title_search VARCHAR(128) NOT NULL DEFAULT '',
    comment_search VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("DATABASE initialization error: %v", err)
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			log.Fatalf("DATABASE initialization error: %v", err)
		}
	}
}

func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("DATABASE close error: %v", err)
		}
	}
}
