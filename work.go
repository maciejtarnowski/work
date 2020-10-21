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
	fmt.Printf(
		"\nToday:\n\tWorked: %s\n\tTotal: %s%s\n",
		result.DayWorked,
		getSignForPositiveDuration(result.DayTotal),
		result.DayTotal,
	)
}

func (a App) Log(weekOffset int) {
	now := time.Now()

	from, to := getDateRangeForLog(now, weekOffset)

	stats, err := statsWork(a.db, from, to)

	if err != nil {
		fmt.Println(err)
		return
	}

	total := stats.Total

	fmt.Printf("%s - %s\n\n", from.Format(dateShortLayout), to.Format(dateShortLayout))
	fmt.Printf(
		"Expected: %s\nWorked: %s\nTotal: %s%s\n",
		stats.Expected,
		stats.Worked,
		getSignForPositiveDuration(total),
		total,
	)

	if len(stats.DayStats) > 0 {
		fmt.Print("\nBy day:\n")

		for _, day := range stats.DayStats {
			dayTotal := day.Total
			fmt.Printf("\t%s: %s%s\n", day.Date.Format(dateShortLayout), getSignForPositiveDuration(dayTotal), dayTotal)
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
