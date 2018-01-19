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
	urlParam := flag.String("url", "", "The URL to watch")
	flag.Parse()
	var url = *urlParam

	if len(url) == 0 {
		fmt.Println("Missing --url")
		return
	}

	// Create the watch
	var modifiedWatch = watches.NewURLModifiedWatch(url)
	if modifiedWatch == nil {
		fmt.Println("Cannot watch url", url)
		return
	}

	var watchMan = watchers.NewWatchMan([]core.Watch{modifiedWatch})

	// Create the triggers
	var modifiedTrigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Println(url, "has been modified at",
			e.Get(watches.URLModifiedTime))
	})

	// Start watching
	fmt.Println("Watching", url)
	watchMan.Watch(modifiedTrigger)

	select {}
}
