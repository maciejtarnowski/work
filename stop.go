package work

import (
	"context"
	"database/sql"
	"time"
)

type StopStats struct {
	SessionWorked time.Duration
	DayWorked     time.Duration
	DayTotal      time.Duration
}

func stopWork(db *sql.DB, when time.Time) (StopStats, error) {
	unfinished, err := checkForUnfinished(db)
	if err != nil {
		return StopStats{}, err
	}
	if !unfinished {
		return StopStats{}, ErrFinishedWork
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(
		ctx,
		`UPDATE work_log SET ended_at = $1 WHERE ended_at IS NULL`,
		when.Unix(),
	)
	if err != nil {
		return StopStats{}, err
	}

	return getStopStats(db, when)
}

func getStopStats(db *sql.DB, when time.Time) (StopStats, error) {
	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()

	var lastSessionSeconds int64
	err := db.QueryRowContext(
		ctx1,
		`SELECT
			ended_at - started_at
		FROM work_log
		WHERE
			started_at IS NOT NULL
			AND ended_at IS NOT NULL
		ORDER BY id DESC
		LIMIT 1`,
	).Scan(&lastSessionSeconds)

	if err != nil {
		return StopStats{}, err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	from := time.Date(when.Year(), when.Month(), when.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(when.Year(), when.Month(), when.Day(), 23, 59, 59, 0, time.Local)
	var daySeconds int64
	err = db.QueryRowContext(
		ctx2,
		`SELECT
			SUM(ended_at - started_at)
		FROM work_log
		WHERE
			started_at >= $1
			AND ended_at <= $2
			AND started_at IS NOT NULL
			AND ended_at IS NOT NULL
		`,
		from.Unix(),
		to.Unix(),
	).Scan(&daySeconds)

	if err != nil {
		return StopStats{}, err
	}

	return StopStats{
		SessionWorked: time.Duration(lastSessionSeconds) * time.Second,
		DayWorked:     time.Duration(daySeconds) * time.Second,
		DayTotal:      (time.Duration(daySeconds) * time.Second) - (8 * time.Hour),
	}, nil
}
