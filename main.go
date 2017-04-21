package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry/gosigar"
)

type Status int

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
	used := int(float64(swap.Used) / float64(swap.Total) * 100)
	fmt.Fprintf(os.Stdout, "Swap usage: %d%%\n", used)

	if used >= f {
		return fail
	}

	if used >= w {
		return warn
	}

	return ok
}

func load(w, f int) Status {
	concreteSigar := sigar.ConcreteSigar{}
	avg, err := concreteSigar.GetLoadAverage()
	if err != nil {
		fmt.Printf("Failed to get load average")
		return ok
	}

	cpus := sigar.CpuList{}
	cpus.Get()

	one := (avg.One / float64(len(cpus.List)))
	five := (avg.Five / float64(len(cpus.List)))
	fifteen := (avg.Fifteen / float64(len(cpus.List)))

	fmt.Printf("Load Average / CPUs: %f %f %f\n", one, five, fifteen)

	if five >= float64(f)/100.0 {
		return fail
	}

	if five >= float64(w)/100.0 {
		return warn
	}

	return ok
}

func mem(w, f int) Status {
	mem := sigar.Mem{}
	mem.Get()
	used := int(float64(mem.ActualUsed) / float64(mem.Total) * 100)
	fmt.Fprintf(os.Stdout, "Memory usage: %d%%\n", used)

	if used >= f {
		return fail
	}

	if used >= w {
		return warn
	}

	return ok
}

const output_format = "%-15s %4s %4s %5s %4s %-15s\n"

func formatSize(size uint64) string {
	return sigar.FormatSize(size * 1024)
}

func disk(w, f int) Status {
	s := ok

	fslist := sigar.FileSystemList{}
	fslist.Get()
	fmt.Fprintf(os.Stdout, output_format,
		"Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	for _, fs := range fslist.List {
		dir_name := fs.DirName

		usage := sigar.FileSystemUsage{}

		usage.Get(dir_name)

		fmt.Fprintf(os.Stdout, output_format,
			fs.DevName,
			formatSize(usage.Total),
			formatSize(usage.Used),
			formatSize(usage.Avail),
			sigar.FormatPercent(usage.UsePercent()),
			dir_name)

		if int(usage.UsePercent()) >= w && s < warn {
			s = warn
		}

		if int(usage.UsePercent()) >= f {
			s = fail
		}
	}

	return s
}
