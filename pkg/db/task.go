package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	if DB == nil {
		return 0, fmt.Errorf("database is not initialized")
	}

	query := `INSERT INTO scheduler(date,title,comment,repeat) VALUES (?,?,?,?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTask(id string) (*Task, error) {
	var task = &Task{}

	if DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

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
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	result, err := DB.Exec(
		`UPDATE scheduler 
		 SET date = ?, title = ?, comment = ?, repeat = ?
		 WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
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

func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

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

	var tasks []*Task

	if DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat
		 FROM scheduler
		 ORDER BY date
		 LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("error with select tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("scanning error: %w", err)
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("scanning error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
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
