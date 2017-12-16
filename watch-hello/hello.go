package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/deanydean/watchit/watches"
	"github.com/deanydean/watchit/watchman"
	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/core/triggers"
)

// Function that will create a new event each time it is called
func newEvent() *core.Event {
	data := map[string]string{
		"Hello": "world!",
	}

	return core.NewEvent(data)
}

func main() {
	// Get cli flags
	pauseFor := flag.Int("pause", 1, "For seconds between hellos")
	watchFor := flag.Int("for", 10, "How long to watch for")
	flag.Parse()

	// A poller watch that will call the newEvent method at intervals
	var poller = watches.NewPollerWatch(newEvent, *pauseFor)

	// The action that will print the new event
	var action = triggers.FuncTrigger(func(e *core.WatchEvent) {
		fmt.Println("Hello,", e.Get("Hello"))
	})

	// The watchman that tells the poller to trigger the action on events
	var watchman = watchman.NewWatchMan(poller, action)

	// Start watching
	fmt.Println("Watching hellos....")
	var canceller = watchman.Watch()

	// Wait for a bit
	time.Sleep(time.Duration(*watchFor) * time.Second)

	// Caller the canceller to cancel the watchman
	canceller()

	fmt.Println("Complete!")
}
