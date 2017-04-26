package siren

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func NTP(w, f int) Status {
	out, err := exec.Command("ntpq", "-c", "rv 0 offset").Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "\n%s: unable to determine NTP drift. %s\n", warn, err)
		return warn
	}

	// Output: offset=72.062
	floatStr := strings.Trim(string(out[7:]), "\n")
	drift, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\n%s: unable to determine NTP drift\n", warn)
		return warn
	}

	fmt.Fprintf(os.Stdout, "NTP drift: %0.2fms\n", drift)

	drift = math.Abs(drift)

	if drift >= float64(f) {
		fmt.Fprintf(os.Stdout, "\n%s: NTP drift exceeds threshold (%0.2fms >= %dms)\n", fail, drift, f)
		return fail
	}

	if drift >= float64(w) {
		fmt.Fprintf(os.Stdout, "\n%s: NTP drift exceeds threshold (%0.2fms >= %dms)\n", warn, drift, w)
		return warn
	}

	return ok
}
