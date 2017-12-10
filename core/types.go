// Package core proces the core watchit API
package core

import "strconv"

// WatchEvent triggered while watching
type WatchEvent struct {
	Data map[string]interface{}
}

// Get the named parameter from this Event
func (e WatchEvent) Get(name string) interface{} {
	return e.Data[name]
}

// GetAsInteger gets the named parameter from this Event as an int
func (e WatchEvent) GetAsInteger(name string) (int, error) {
	return strconv.Atoi(e.Data[name].(string))
}

// GetAsFloat gets the named parameter from this Event as a float
func (e WatchEvent) GetAsFloat(name string) (float64, error) {
	return strconv.ParseFloat(e.Data[name].(string), 64)
}

// NewWatchEvent creates a new WatchEvent with the provided event data
func NewWatchEvent(data map[string]interface{}) *WatchEvent {
	event := new(WatchEvent)
	event.Data = data
	return event
}

// WatchTrigger is a type that can be triggered on a watch event
type WatchTrigger interface {
	OnEvent(*WatchEvent)
}

// WatcherCanceller can cancel a Watcher
type WatcherCanceller func()

// Watcher is an interface for things that Watch
type Watcher interface {
	Watch(WatchTrigger) WatcherCanceller
}

// Watch is an interface for something that can be watched
type Watch interface {
	Observe() *WatchEvent
}
