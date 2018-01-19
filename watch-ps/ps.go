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

	// Create process watches for CPU and Mem usage and one to see when the
	// process ends
	var deathWatch = watches.NewProcessDeathWatch(pid)
	var highCPUWatch = watches.NewProcessHighCPUWatch(pid, 50)
	var highMemWatch = watches.NewProcessHighMemWatch(pid, 500000000)

	var watchMan = watchers.NewWatchMan([]core.Watch{highCPUWatch, highMemWatch})

	// Create the triggers
	var textOutputTrigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Println("Process", pid, e)
	})
	// Start watching
	fmt.Println("Watching pid", pid, "....")
	var cancelWatch = watchMan.Watch(textOutputTrigger)
	processEnded := make(chan bool)

	go func() {
		// Action for when process dies
		var onDeath = func(e *core.WatchEvent) {
			fmt.Println("Process", pid, "died!")

			// Caller the cancellers to cancel the watchman
			cancelWatch()

			processEnded <- true
		}
		var deathWatchMan = watchers.NewWatchMan([]core.Watch{deathWatch})
		deathWatchMan.Watch(triggers.NewFuncTrigger(onDeath))
	}()

	// Need to wait for the death watcher....
	<-processEnded
}
