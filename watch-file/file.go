package main

import (
	"flag"
	"fmt"

	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/core/triggers"
	"github.com/oddcyborg/watchit/core/watchers"
	"github.com/oddcyborg/watchit/core/watches"
)

func main() {
	// Get cli flags
	fileParam := flag.String("file", "", "The name of the file to watch")
	flag.Parse()

	var fileName = *fileParam
	if len(fileName) == 0 {
		fmt.Println("Missing --file")
		return
	}

	var modifiedWatch = watches.NewFileModifiedWatch(fileName)
	var watchMan = watchers.NewWatchMan(modifiedWatch)

	// Create the triggers
	var modifiedTrigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Println(fileName, "has been modified at",
			e.Get(watches.FileModTime))
	})

	// Start watching
	fmt.Println("Watching", fileName)
	watchMan.Watch(modifiedTrigger)

	select {}
}
