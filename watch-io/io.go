package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/core/triggers"
	"github.com/oddcyborg/watchit/core/utils"
	"github.com/oddcyborg/watchit/core/watchers"
	"github.com/oddcyborg/watchit/core/watches"
)

func main() {
	// Get cli flags
	pidFlag := flag.Int("pid", -1, "The pid ")
	flag.Parse()

	var pid = *pidFlag
	if pid <= 0 {
		fmt.Fprintln(os.Stderr, "Invalid pid", pid)
		os.Exit(1)
	} else if !utils.ProcessExists(pid) {
		fmt.Fprintln(os.Stderr, "No such process with pid", pid)
		os.Exit(1)
	}

	var ioWatch = watches.NewProcessHighIOWatch(pid, 10)
	var trigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Println("Process", pid, "has high io:")
	})
	var watchMan = watchers.NewWatchMan(ioWatch)
	watchMan.Watch(trigger)
	select {}
}
