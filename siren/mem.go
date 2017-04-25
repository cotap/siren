package siren

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/gosigar"
)

func Mem(w, f int) Status {
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
