package watches

import "github.com/deanydean/watchit/core"
import "time"

// A Watch that will poll an EventSupplier for events at regular intervals
type PollerWatch struct {
    poll core.EventSupplier
    interval int
    polling bool
}

// Start the PollerWatch polling for events
func (pw PollerWatch) Start(handler core.EventConsumer) core.WatchCanceller {
    pw.polling = true

    go func() {
        for pw.polling {
            result := pw.poll()
            if result != nil {
                handler(result)
            }
            time.Sleep( time.Duration(pw.interval) * time.Second )
        }
    }()

    // Return an activewatch that will end polling when cancel() is called
    return func() { pw.polling = false; }
}

// Create a new PollerWatch that will poll at the provided interval
func NewPollerWatch(poll core.EventSupplier, interval int) *PollerWatch {
    poller := new(PollerWatch)
    poller.poll = poll
    poller.interval = interval
    poller.polling = false
    return poller
}
