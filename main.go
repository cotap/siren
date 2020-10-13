package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cotap/siren/siren"
)

const usage = "siren [mem|swap|disk|load|chrony|proc] ((WARN_LEVEL FAIL_LEVEL) | PID)"

func main() {
	cmd, w, f, pid, err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch cmd {
	case "mem":
		os.Exit(int(siren.Mem(w, f)))
	case "swap":
		os.Exit(int(siren.Swap(w, f)))
	case "disk":
		os.Exit(int(siren.Disk(w, f)))
	case "load":
		os.Exit(int(siren.Load(w, f)))
	case "chrony":
		os.Exit(int(siren.CHRONY(w, f)))
	case "proc":
		os.Exit(int(siren.Proc(pid)))
	default:
		fmt.Println(usage)
	}
}

func parseArgs() (string, int, int, int, error) {
	var (
		cmd       string
		w, f, pid int
		err       error
	)

	if len(os.Args) < 3 {
		return cmd, w, f, pid, errors.New(usage)
	}

	cmd = os.Args[1]

	switch cmd {
	case "mem", "swap", "disk", "load", "ntp":
		if len(os.Args) < 4 {
			return cmd, w, f, pid, errors.New(usage)
		}

		w, err = strconv.Atoi(os.Args[2])
		if err != nil {
			err = errors.New("WARN_LEVEL must be an integer")
		}

		f, err = strconv.Atoi(os.Args[3])
		if err != nil {
			err = errors.New("FAIL_LEVEL must be an integer")
		}
	case "proc":
		pid, err = strconv.Atoi(os.Args[2])
	}

	return cmd, w, f, pid, err
}
