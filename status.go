package work

import (
	"context"
	"database/sql"
	"time"
)

type Status struct {
	IsActive  bool
	Worked    time.Duration
	Remaining time.Duration
}

func workStatus(db *sql.DB, when time.Time) (Status, error) {
	unfinished, err := checkForUnfinished(db)
	if err != nil {
		return Status{}, err
	}

	dayWorked, err := getWorkTimeForDay(db, when)
	if err != nil {
		return Status{}, err
	}

	from := time.Date(when.Year(), when.Month(), when.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(when.Year(), when.Month(), when.Day(), 23, 59, 59, 0, time.Local)

	return Status{
		IsActive:  unfinished,
		Worked:    dayWorked,
		Remaining: calculateExpected(from, to) - dayWorked,
	}, nil
}

func getWorkTimeForDay(db *sql.DB, when time.Time) (time.Duration, error) {
	var seconds sql.NullInt64

	from := time.Date(when.Year(), when.Month(), when.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(when.Year(), when.Month(), when.Day(), 23, 59, 59, 0, time.Local)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.QueryRowContext(
		ctx,
		`SELECT
			SUM(IFNULL(ended_at, $1) - started_at)
		FROM work_log
		WHERE
			started_at >= $2
			AND (ended_at <= $3 OR ended_at IS NULL)
			AND started_at IS NOT NULL`,
		when.Unix(),
		from.Unix(),
		to.Unix(),
	).Scan(&seconds)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	if seconds.Valid {
		return time.Duration(seconds.Int64) * time.Second, nil
	}
	return 0, nil
}
