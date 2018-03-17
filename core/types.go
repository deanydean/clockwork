// Types for the core clockwork API
package core

import (
	"strconv"
	"time"
)

// WatchEvent triggered while watching
type WatchEvent struct {
	code      int
	stopWatch bool
	timestamp time.Time
	Data      map[string]interface{}
}

// Get the named parameter from this Event
func (e WatchEvent) Get(name string) interface{} {
	return e.Data[name]
}

func (e WatchEvent) GetAsString(name string) string {
	return e.Data[name].(string)
}

// GetAsInteger gets the named parameter from this Event as an int
func (e WatchEvent) GetAsInteger(name string) int {
	var value, err = strconv.Atoi(e.Data[name].(string))
	if err != nil {
		return -1
	}
	return value
}

// GetAsFloat gets the named parameter from this Event as a float
func (e WatchEvent) GetAsFloat(name string) float64 {
	var value, err = strconv.ParseFloat(e.Data[name].(string), 64)
	if err != nil {
		return -1
	}
	return value
}

func (e WatchEvent) GetAsArray(name string) []interface{} {
	return e.Data[name].([]interface{})
}

func (e WatchEvent) ShouldStop() bool {
	return e.stopWatch
}

func (e WatchEvent) GetTime() time.Time {
	return e.timestamp
}

func (e WatchEvent) Result() int {
	return e.code
}

func (e *WatchEvent) SetStatus(code int, stop bool) {
	e.code = code
	e.stopWatch = stop
}

// NewWatchEvent creates a new WatchEvent with the provided event data
func NewWatchEvent(data map[string]interface{}) *WatchEvent {
	event := new(WatchEvent)
	event.timestamp = time.Now()
	event.code = 0
	event.stopWatch = false
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
