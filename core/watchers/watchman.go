package watchers

import (
	"time"

	"github.com/oddcyborg/watchit/core"
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
		var polling = true
		for polling {
			for watch := range wm.watches {
				result := wm.watches[watch].Observe()
				if result != nil {
					log.Debug("Got result=%s from watch=%s", result.Data,
						wm.watches[watch])
					trigger.OnEvent(result)
				}
			}

			select {
			case stopSignal := <-wm.stopper:
				polling = !stopSignal
			default:
				time.Sleep(time.Duration(wm.interval) * time.Second)
			}
		}
	}()

	// Return an activewatch that will end polling when cancel() is called
	return wm.Stop
}

// Stop the PollerWatch polling for events
func (wm WatchMan) Stop() {
	wm.stopper <- true
}
