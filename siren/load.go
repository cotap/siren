package siren

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/gosigar"
)

func Load(w, f int) Status {
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
