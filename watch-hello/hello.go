package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/deanydean/watchit/core/watchers"
	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/core/triggers"
)

// HelloWatch will create a new event each time it is observed
type HelloWatch struct {
}

// Observe the hello events
func (hw HelloWatch) Observe() *core.WatchEvent {
	data := map[string]interface{}{
		"Hello": "world!",
	}

	return core.NewWatchEvent(data)
}

func main() {
	// Get cli flags
	pauseFor := flag.Int("pause", 1, "For seconds between hellos")
	watchFor := flag.Int("for", 10, "How long to watch for")
	flag.Parse()

	// A poller watcher that will call the newEvent method at intervals
	var poller = watchers.NewPollerWatcher(new(HelloWatch), *pauseFor)

	// The action that will print the new event
	var action = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Println("Hello,", e.Get("Hello"))
	})

	// Start watching
	var canceller = poller.Watch(action)
	fmt.Println("Watching hellos....")

	// Wait for a bit
	time.Sleep(time.Duration(*watchFor) * time.Second)

	// Caller the canceller to cancel the watchman
	canceller()

	fmt.Println("Complete!")
}
