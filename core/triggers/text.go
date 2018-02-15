package triggers

import (
	"fmt"

	"github.com/deanydean/clockwork/core"
)

// TextReporterTrigger reports using text when a watch event triggers
type TextReporterTrigger struct {
	message string
}

// OnEvent is called when a watch event triggers
func (trigger TextReporterTrigger) OnEvent(event *core.WatchEvent) {
	fmt.Printf(trigger.message, event)
}

// NewTextReporterTrigger create a new TextReporterTrigger
func NewTextReporterTrigger(message string) TextReporterTrigger {
	var trigger = new(TextReporterTrigger)
	trigger.message = message
	return *trigger
}
