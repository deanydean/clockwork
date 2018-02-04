package watchers

import (
    "time"

    "github.com/deanydean/clockwork/core"
    "github.com/deanydean/clockwork/core/utils"
)

var log = utils.GetLogger()

// PollerWatcher will poll an EventSupplier for events at regular intervals
type PollerWatcher struct {
    poll     core.Watch
    interval int
    stopper  chan bool
}

// Watch the Watch by polling it's Observe function
func (pw PollerWatcher) Watch(trigger core.WatchTrigger) core.WatcherCanceller {
    go func() {
        var polling = true
        for polling {
            result := pw.poll.Observe()
            if result != nil {
                log.Debug("Got result=%s from watch=%s", result.Data, pw.poll)
                trigger.OnEvent(result)
            }

            select {
            case stopSignal := <-pw.stopper:
                polling = !stopSignal
            default:
                time.Sleep(time.Duration(pw.interval) * time.Second)
            }
        }
    }()

    // Return an activewatch that will end polling when cancel() is called
    return pw.Stop
}

// Stop the PollerWatch polling for events
func (pw PollerWatcher) Stop() {
    pw.stopper <- true
}

// NewPollerWatcher creates a new PollerWatch that will poll at the provided
// interval (in seconds)
func NewPollerWatcher(poll core.Watch, interval int) *PollerWatcher {
    poller := new(PollerWatcher)
    poller.poll = poll
    poller.interval = interval
    poller.stopper = make(chan bool, 1)
    return poller
}
