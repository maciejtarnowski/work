package work

import (
	"database/sql"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02 15:04:05"
const dateShortLayout = "2006-01-02"

type App struct {
	db *sql.DB
}

func NewApp(dbPath string) App {
	db, err := connectToDB(dbPath)
	if err != nil {
		panic(err)
	}

	return App{db: db}
}

func (a App) Close() error {
	return a.db.Close()
}

func (a App) Start() {
	now := time.Now()

	err := startWork(a.db, now)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Work started at %s\n", now.Format(dateLayout))
}

func (a App) Stop() {
	now := time.Now()

	result, err := stopWork(a.db, now)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Work stopped at %s (after %s)\n", now.Format(dateLayout), result.SessionWorked)
	fmt.Printf("\nToday:\n\tWorked: %s\n\tTotal: %s\n", result.DayWorked, result.DayTotal)
}

func (a App) Log() {
	now := time.Now()

	dayOffset := now.Weekday() - 1
	if dayOffset == -1 {
		dayOffset = 6
	}
	from := time.Date(now.Year(), now.Month(), now.Day()-int(dayOffset), 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	stats, err := statsWork(a.db, from, to)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s - %s\n\n", from.Format(dateShortLayout), to.Format(dateShortLayout))
	fmt.Printf(
		"Expected: %s\nWorked: %s\nTotal: %s\n",
		stats.Expected.Round(time.Minute),
		stats.Worked.Round(time.Minute),
		stats.Total.Round(time.Minute),
	)

	if len(stats.DayStats) > 0 {
		fmt.Print("\nBy day:\n")

		for _, day := range stats.DayStats {
			fmt.Printf("\t%s: %s\n", day.Date.Format(dateShortLayout), day.Total.Round(time.Minute))
		}
	}
}

func (a App) Status() {
	now := time.Now()

	status, err := workStatus(a.db, now)
	if err != nil {
		fmt.Println(err)
		return
	}

	if status.IsActive {
		fmt.Println("Work: IN PROGRESS")
	} else {
		fmt.Println("Work: STOPPED")
	}

	fmt.Printf("\nToday:\n\tWorked: %s\n\tRemaining: %s\n", status.Worked, status.Remaining)
}
