package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/noonien/pm"
)

var (
	timeDelay  = flag.String("t", "now", "time after which to execute operation")
	background = flag.Bool("bg", false, "run in background when a time is specified")
)

func main() {
	ops := []struct {
		name  string
		check func() bool
		fn    func() error
	}{
		{"poweroff", powermanager.CanPowerOff, powermanager.PowerOff},
		{"reboot", powermanager.CanReboot, powermanager.Reboot},
		{"suspend", powermanager.CanSuspend, powermanager.Suspend},
		{"hybrid-sleep", powermanager.CanHybridSleep, powermanager.HybridSleep},
		{"hibernate", powermanager.CanHibernate, powermanager.Hibernate},
	}

	flag.Usage = func() {
		var availableOps []string
		for _, op := range ops {
			if op.check() {
				availableOps = append(availableOps, op.name)
			}
		}
		if len(availableOps) == 0 {
			fmt.Fprintln(os.Stderr, "Sorry, no available operations for your computer")
			os.Exit(1)
		}

		fmt.Printf("Usage: %s [OPTION]... [%s]\n", os.Args[0], strings.Join(availableOps, "|"))
		fmt.Println("User space power management")
		fmt.Println()

		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	op := flag.Arg(0)
	var opFn func() error
	for _, o := range ops {
		if o.name == op {
			if o.check() {
				opFn = o.fn
			}
			break
		}
	}
	if opFn == nil {
		flag.Usage()
		fmt.Println()

		fmt.Fprintln(os.Stderr, "Invalid operation:", op)
		os.Exit(1)
	}

	if *timeDelay != "now" {
		duration, err := time.ParseDuration(*timeDelay)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Couldn't parse time:", err)
			os.Exit(1)
		}

		if *background {
			cmd := exec.Command(os.Args[0], "-t", *timeDelay, op)
			err = cmd.Start()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Couldn't start in background:", err)
				os.Exit(1)
			}
			fmt.Printf("Will enter %s state after %v, run `kill %d` to cancel\n", op, duration, cmd.Process.Pid)
			os.Exit(0)
		}

		fmt.Printf("Will enter %s state after %v\n", op, duration)
		time.Sleep(duration)
	}

	err := opFn()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
