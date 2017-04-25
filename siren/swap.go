package siren

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/gosigar"
)

func Swap(w, f int) Status {
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
