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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = db.ExecContext(ctx, `INSERT INTO work_log(started_at) VALUES ($1)`, when.Unix())

	return err
}
