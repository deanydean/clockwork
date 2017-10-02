package watchman

import "github.com/oddcyborg/watchit/core"

// A simple Watcher that links a Watch to an EventConsumer
type WatchMan struct {
    watch core.Watch
    onEvent core.EventConsumer
}

// Create a new WatchMan
func NewWatchMan(watch core.Watch, onEvent core.EventConsumer) *WatchMan {
    wm := new(WatchMan)
    wm.watch = watch
    wm.onEvent = onEvent
    return wm
}

// Start watching
func (wm WatchMan) Watch() core.WatchCanceller {
    return wm.watch.Start(wm.onEvent)
}
