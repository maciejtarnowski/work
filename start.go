package work

import (
	"context"
	"database/sql"
	"time"
)

func startWork(db *sql.DB, when time.Time) error {
	unfinished, err := checkForUnfinished(db)
	if err != nil {
		return err
	}
	if unfinished {
		return ErrUnfinishedWork
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, `INSERT INTO work_log(started_at) VALUES ($1)`, when.Unix())

	return err
}
