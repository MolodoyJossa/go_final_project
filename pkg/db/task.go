package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
)

func toSearch(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type Task struct {
	ID           string `json:"id"`
	Date         string `json:"date"`
	Title        string `json:"title"`
	Comment      string `json:"comment"`
	Repeat       string `json:"repeat"`
	TitleLower   string `json:"title_search"`
	CommentLower string `json:"comment_search"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler(date,title,comment,repeat,title_search,comment_search) VALUES (?,?,?,?,?,?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, toSearch(task.Title), toSearch(task.Comment))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func GetTask(id string) (*Task, error) {
	var task = &Task{}

	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat
		 FROM scheduler
		 WHERE id = ?`,
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error with get task: %w", err)
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	result, err := DB.Exec(
		`UPDATE scheduler 
		 SET date = ?, title = ?, comment = ?, repeat = ?, title_search = ?, comment_search = ?
		 WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, toSearch(task.Title), toSearch(task.Comment), task.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateTaskDate(id string, nextDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, nextDate, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	result, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func Tasks(limit int) ([]*Task, error) {

	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat
		 FROM scheduler
		 ORDER BY date
		 LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("error with select tasks: %w", err)
	}
	defer rows.Close()

	return rowsScan(rows)
}

func SearchTasks(search string, limit int) ([]*Task, error) {
	if searchDate, err := time.Parse("02.01.2006", search); err == nil {
		rows, err := DB.Query(`SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE date = ? 
	          ORDER BY id
	          LIMIT ?`, searchDate.Format("20060102"), limit)
		if err != nil {
			return nil, fmt.Errorf("failed to query tasks by date: %w", err)
		}
		defer rows.Close()
		return rowsScan(rows)
	} else {
		search = "%" + toSearch(search) + "%"
		rows, err := DB.Query(`SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE title_search LIKE ? OR comment_search LIKE ?
	          ORDER BY date
	          LIMIT ?`, search, search, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to query tasks by search: %v", err)
		}
		defer rows.Close()
		return rowsScan(rows)
	}
}

func rowsScan(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("scanning error: %w", err)
		}
		tasks = append(tasks, task)
	}
	err := rows.Err()
	if err != nil {
		return nil, fmt.Errorf("scanning error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
