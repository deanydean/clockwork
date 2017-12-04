package watchers

import "github.com/oddcyborg/watchit/core"

// WatchMan is a Watcher that links a Watch to an EventConsumer
type WatchMan struct {
	watch  core.Watch
	poller *PollerWatcher
}

// NewWatchMan creates a new WatchMan
func NewWatchMan(watch core.Watch) *WatchMan {
	wm := new(WatchMan)
	wm.watch = watch
	wm.poller = NewPollerWatcher(watch, 1)

	return wm
}

// Watch tells the WatchMan to start watching
func (wm WatchMan) Watch(trigger core.WatchTrigger) core.WatcherCanceller {
	return wm.poller.Watch(trigger)
}
