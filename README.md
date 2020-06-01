# work

Simple work time log

## Overview

Keep track of how long you've worked each day.

 - Written in Go, compiles to a single binary
 - Uses SQLite as persistent storage

## Installation

Requires Go 1.14

```
$ go get -u github.com/maciejtarnowski/work
$ go install github.com/maciejtarnowski/work/cmd/work
```

Make sure you have `$GOPATH/bin` in your `$PATH`:
```
export PATH="$PATH:$GOPATH/bin"
```

## Usage

**Warning:** Running any of the commands automatically creates `work.db` file inside your `$HOME` directory.

### Start work

Simply run:
```
$ work start
```

### Pause or finish work

Run:
```
$ work stop
```

Stats are grouped by day, so it doesn't matter how many times you start and stop your work each day.

### Display log

Run:
```
$ work log
```

The command prints work log for the current week (starting on Monday), for example:
```
2020-05-18 - 2020-05-23

Expected: 40h0m0s
Worked: 8m0s
Total: -39h52m0s

By day:
	2020-05-23: -7h52m0s
```

It currently assumes that work day equals 8 hours and skips Saturdays and Sundays.

Only finished work sessions are included in the log.

### Work status

Run:
```
$ work status
```
or shorter:
```
$ work st
```

This command prints the current status of your work, for example:
```
Work: IN PROGRESS

Today:
	Worked: 30m45s
	Remaining: 7h29m15s
```

Stats include the ongoing work session.

## State of the project

It's in prototype/MVP stage. I want to find out how it works for me and decide if I should spend more time developing it.

## License

MIT
