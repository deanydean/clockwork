package triggers

import (
	"fmt"

	"github.com/oddcyborg/watchit/core"
)

// TextReporterTrigger
type TextReporterTrigger struct {
	message string
}

// OnEvent
func (trigger TextReporterTrigger) OnEvent(event core.WatchEvent) {
	fmt.Printf(trigger.message, event)
}
