package work

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var (
	ErrUnfinishedWork = errors.New("cannot start work when there is unfinished one")
	ErrFinishedWork   = errors.New("cannot stop work when there is no unfinished one")
)

func checkForUnfinished(db *sql.DB) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := db.QueryRowContext(ctx, `SELECT TRUE FROM work_log WHERE ended_at IS NULL LIMIT 1`).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return exists, nil
}

func connectToDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	createTablesStmt := `
		CREATE TABLE IF NOT EXISTS work_log(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			started_at INTEGER NULL,
			ended_at INTEGER NULL
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, createTablesStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getSignForPositiveDuration(d time.Duration) string {
	if d > 0 {
		return "+"
	}
	return ""
}
