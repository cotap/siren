package siren

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/gosigar"
)

func Proc(pid int) Status {
	list := sigar.ProcList{}
	list.Get()

	for p := range list.List {
		if p == pid {
			fmt.Fprintf(os.Stdout, "\n%s: process with PID %d is running\n", ok, pid)
			return ok
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s: process with PID %d is NOT running\n", fail, pid)
	return fail
}
