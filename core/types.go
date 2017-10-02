// The core watchit API types
package core

// Event triggered while watching
type Event struct {
    data map[string]string
}

// Get the named parameter from this Event
func (e Event) Get(name string) string {
    return e.data[name]
}

// Create a new Event with the provided event data
func NewEvent(data map[string]string) *Event {
    event := new(Event)
    event.data = data
    return event
}

// A supplier of events
type EventSupplier func() *Event

// A consumer of events
type EventConsumer func(*Event)

// A function that can cancel a Watch
type WatchCanceller func()

// An interface for things that Watch
type Watcher interface {
    Watch() WatchCanceller
}

// An interface for a Watch
type Watch interface {
    Start(EventConsumer) WatchCanceller
}
