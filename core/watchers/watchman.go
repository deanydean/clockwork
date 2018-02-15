package watchers

import (
	"time"

	"github.com/deanydean/clockwork/core"
)

// WatchMan is a Watcher that links a number of Watches to a WatchTrigger
type WatchMan struct {
	watches  []core.Watch
	trigger  *core.WatchTrigger
	interval int
	stopper  chan bool
}

// NewWatchMan creates a new WatchMan
func NewWatchMan(watches []core.Watch) *WatchMan {
	wm := new(WatchMan)
	wm.watches = watches

	wm.interval = 1
	wm.stopper = make(chan bool, 1)

	return wm
}

// Watch tells the WatchMan to start watching
func (wm WatchMan) Watch(trigger core.WatchTrigger) core.WatcherCanceller {
	go func() {
		// Poll until polling is unset
		var polling = true
		for polling {
			for w := range wm.watches {
				watch := wm.watches[w]

				// Observe the watch value
				go func() {
					result := watch.Observe()
					if result != nil {
						log.Debug("Got result=%s from watch=%s", result.Data,
							watch)

						// Send the trigger
						go func() {
							trigger.OnEvent(result)
						}()
					}
				}()
			}

			select {
			case stopSignal := <-wm.stopper:
				polling = !stopSignal
			default:
				time.Sleep(time.Duration(wm.interval) * time.Second)
			}
		}
	}()

	// Return an WatcherCanceller that will end polling when called
	return wm.Stop
}

// Stop watching for events
func (wm WatchMan) Stop() {
	wm.stopper <- true
}
