package work

import (
	"context"
	"database/sql"
	"time"
)

func stopWork(db *sql.DB, when time.Time) error {
	unfinished, err := checkForUnfinished(db)
	if err != nil {
		return err
	}
	if !unfinished {
		return ErrFinishedWork
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = db.ExecContext(
		ctx,
		`UPDATE work_log SET ended_at = $1 WHERE ended_at IS NULL`,
		when.Unix(),
	)

	return err
}
