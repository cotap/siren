package siren

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/gosigar"
)

const diskFormat = "%-10s %-15s %4s %4s %5s %4s %-15s\n"

func formatSize(size uint64) string {
	return sigar.FormatSize(size * 1024)
}

func Disk(w, f int) Status {
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
