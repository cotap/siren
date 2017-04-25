package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry/gosigar"
)

type Status int

func (s Status) String() string {
	switch s {
	case ok:
		return "OK"
	case warn:
		return "WARN"
	default:
		return "FAIL"
	}
}

const (
	ok Status = iota
	warn
	fail
)

func main() {

	if len(os.Args) != 4 {
		fmt.Println("siren [mem|swap|disk|load] WARN_LEVEL FAIL_LEVEL")
		return
	}

	cmd := os.Args[1]

	w, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("WARN_LEVEL must be an integer")
	}

	f, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("FAIL_LEVEL must be an integer")
	}

	switch cmd {
	case "mem":
		os.Exit(int(mem(w, f)))
	case "swap":
		os.Exit(int(swap(w, f)))
	case "disk":
		os.Exit(int(disk(w, f)))
	case "load":
		os.Exit(int(load(w, f)))
	default:
		fmt.Println("One of the following commands: mem, swap, disk, load\n")
	}
}

func swap(w, f int) Status {
	swap := sigar.Swap{}
	swap.Get()
	used := float64(swap.Used) / float64(swap.Total) * 100.0

	if used <= 0 {
		fmt.Fprintf(os.Stdout, "No swap usage")
		return ok
	}

	fmt.Fprintf(os.Stdout, "Swap usage: %0.2f%%\n", used)

	if used >= float64(f) {
		fmt.Fprintf(os.Stdout, "\n%s: swap usage exceeds threshold (%0.2f%% >= %d%%)\n", fail, used, f)
		return fail
	}

	if used >= float64(w) {
		fmt.Fprintf(os.Stdout, "\n%s: swap usage exceeds threshold (%0.2f%% >= %d%%)\n", warn, used, w)
		return warn
	}

	return ok
}

func load(w, f int) Status {
	concreteSigar := sigar.ConcreteSigar{}
	avg, err := concreteSigar.GetLoadAverage()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed to get load average")
		return ok
	}

	cpus := sigar.CpuList{}
	cpus.Get()

	numCPUs := len(cpus.List)

	one := (avg.One / float64(numCPUs)) * 100.0
	five := (avg.Five / float64(numCPUs)) * 100.0
	fifteen := (avg.Fifteen / float64(numCPUs)) * 100.0

	fmt.Fprintf(os.Stdout, "CPUs: %d\n", numCPUs)
	fmt.Fprintf(os.Stdout, "Load Averages: %0.3f %0.3f %0.3f\n", avg.One, avg.Five, avg.Fifteen)
	fmt.Fprintf(os.Stdout, "Normalized Load: %0.2f%% %0.2f%% %0.2f%%\n", one, five, fifteen)

	if five >= float64(f) {
		fmt.Fprintf(os.Stdout, "\n%s: 5min normalized load exceeds threshold (%0.2f%% >= %d%%)\n", fail, five, f)
		return fail
	}

	if five >= float64(w) {
		fmt.Fprintf(os.Stdout, "\n%s: 5min normalized load exceeds threshold (%0.2f%% >= %d%%)\n", warn, five, w)
		return warn
	}

	return ok
}

func mem(w, f int) Status {
	mem := sigar.Mem{}
	mem.Get()
	used := float64(mem.ActualUsed) / float64(mem.Total) * 100
	fmt.Fprintf(os.Stdout, "Memory usage: %0.2f%%\n", used)

	if used >= float64(f) {
		fmt.Fprintf(os.Stdout, "\n%s: memory usage exceeds threshold (%0.2f%% >= %d%%)\n", fail, used, f)
		return fail
	}

	if used >= float64(w) {
		fmt.Fprintf(os.Stdout, "\n%s: memory usage exceeds threshold (%0.2f%% >= %d%%)\n", warn, used, w)
		return warn
	}

	return ok
}

const diskFormat = "%-10s %-15s %4s %4s %5s %4s %-15s\n"

func formatSize(size uint64) string {
	return sigar.FormatSize(size * 1024)
}

func disk(w, f int) Status {
	s := ok

	fslist := sigar.FileSystemList{}
	fslist.Get()
	fmt.Fprintf(os.Stdout, diskFormat,
		"Status", "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	for _, fs := range fslist.List {
		status := ok
		dirDame := fs.DirName
		usage := sigar.FileSystemUsage{}
		usage.Get(dirDame)

		if usage.UsePercent() >= float64(w) {
			status = warn
		}

		if usage.UsePercent() >= float64(f) {
			status = fail
		}

		fmt.Fprintf(os.Stdout, diskFormat,
			status,
			fs.DevName,
			formatSize(usage.Total),
			formatSize(usage.Used),
			formatSize(usage.Avail),
			sigar.FormatPercent(usage.UsePercent()),
			dirDame)

		if status > s {
			s = status
		}

	}

	return s
}
