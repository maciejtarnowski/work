package main

import (
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
		fmt.Print("usage:\n\twork start\n\twork stop\n\twork log\n")
		return
	}

	cmd := os.Args[1]

	switch strings.TrimSpace(cmd) {
	case "start":
		app.Start()
	case "stop":
		app.Stop()
	case "log":
		app.Log()
	default:
		fmt.Printf("unknown command: %s\n", strings.TrimSpace(cmd))
	}
}
