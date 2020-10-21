package main

import (
	"flag"
	"fmt"
	"github.com/maciejtarnowski/work"
	"os"
	"strings"
)

func main() {
	dbPath := os.Getenv("HOME") + "/work.db"

	app := work.NewApp(dbPath)
	defer app.Close()

	if len(os.Args) < 2 {
		fmt.Print("usage:\n\twork start\n\twork stop\n\twork log [-week-offset=<n>]\n\twork status or st\n")
		return
	}

	cmd := strings.TrimSpace(os.Args[1])

	switch cmd {
	case "start":
		app.Start()
	case "stop":
		app.Stop()
	case "log":
		logFlags := flag.NewFlagSet("log", flag.ExitOnError)
		weekOffset := logFlags.Int("week-offset", 0, "Display log for a week in the past, e.g. 0 is current week, 1 is last week, etc.")
		logFlags.Parse(os.Args[2:])

		app.Log(*weekOffset)
	case "status":
		app.Status()
	case "st":
		app.Status()
	default:
		fmt.Printf("unknown command: %s\n", cmd)
	}
}
