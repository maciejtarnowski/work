package work

import (
	"context"
	"database/sql"
	"time"
)

type DayStats struct {
	Date time.Time
	Total time.Duration
	Worked time.Duration
	Expected time.Duration
}

type Stats struct {
	Total    time.Duration
	Worked   time.Duration
	Expected time.Duration
	DayStats []DayStats
	From time.Time
	To time.Time
}

func statsWork(db *sql.DB, from, to time.Time) (Stats, error) {
	expected := calculateExpected(from, to)
	stats := Stats{Expected: expected, From: from, To: to}

	total, err := getTotalDuration(db, from, to)
	if err != nil {
		return Stats{}, err
	}
	stats.Worked = total
	stats.Total = total - expected

	dayStats, err := getDayStats(db, from, to)
	if err != nil {
		return Stats{}, nil
	}
	stats.DayStats = dayStats

	return stats, nil
}

func calculateExpected(from time.Time, to time.Time) time.Duration {
	result := 0*time.Hour

	for d := from; d.Day() <= to.Day() && d.Month() <= to.Month() && d.Year() <= to.Year(); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Sunday && d.Weekday() != time.Saturday {
			result += 8*time.Hour
		}
	}

	return result
}

func getTotalDuration(db *sql.DB, from, to time.Time) (time.Duration, error) {
	var seconds sql.NullInt64

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.QueryRowContext(
		ctx,
		`SELECT
			SUM(ended_at - started_at)
		FROM work_log
		WHERE
			started_at >= $1
			AND ended_at <= $2
			AND started_at IS NOT NULL
			AND ended_at IS NOT NULL`,
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
		return time.Duration(seconds.Int64)*time.Second, nil
	}
	return 0, nil
}

func getDayStats(db *sql.DB, from, to time.Time) ([]DayStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(
		ctx,
		`SELECT
			STRFTIME('%Y-%m-%d', started_at, 'unixepoch') AS day,
			SUM(ended_at - started_at) AS duration
		FROM work_log
		WHERE
			started_at >= $1
			AND ended_at <= $2
			AND started_at IS NOT NULL
			AND ended_at IS NOT NULL
		GROUP BY 1`,
		from.Unix(),
		to.Unix(),
	)

	if err != nil {
		return nil, err
	}

	var stats []DayStats

	for rows.Next() {
		var dayStr string
		var seconds int64

		err := rows.Scan(&dayStr, &seconds)

		if err != nil {
			return nil, err
		}

		date, err := time.ParseInLocation("2006-01-02", dayStr, time.Local)
		if err != nil {
			return nil, err
		}

		item := DayStats{}
		item.Date = date
		item.Expected = 8*time.Hour
		item.Worked = time.Duration(seconds)*time.Second
		item.Total = item.Worked - item.Expected

		stats = append(stats, item)
	}

	return stats, nil
}