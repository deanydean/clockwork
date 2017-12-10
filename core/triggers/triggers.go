package triggers

import "github.com/oddcyborg/watchit/core"

// FuncTrigger calls a function when a WatchEvent triggers
type FuncTrigger struct {
	onEvent func(*core.WatchEvent)
}

// OnEvent is called when a WatchEvent triggers
func (trigger FuncTrigger) OnEvent(event *core.WatchEvent) {
	trigger.onEvent(event)
}

// NewFuncTrigger create a new FuncTrigger for the provided func
func NewFuncTrigger(onEvent func(*core.WatchEvent)) *FuncTrigger {
	trigger := new(FuncTrigger)
	trigger.onEvent = onEvent
	return trigger
}

// BroadcastTrigger sends a WatchEvent to all triggers attaches to this trigger
// when a WatchEvent is triggered
type BroadcastTrigger struct {
	triggers []core.WatchTrigger
}

// OnEvent is called when a WatchEvent triggers
func (bt BroadcastTrigger) OnEvent(e *core.WatchEvent) {
	for t := range bt.triggers {
		bt.triggers[t].OnEvent(e)
	}
}

// NewBroadcastTrigger creates a new BroadcastTrigger for the provided triggers
func NewBroadcastTrigger(triggers []core.WatchTrigger) *BroadcastTrigger {
	bt := new(BroadcastTrigger)
	bt.triggers = triggers
	return bt
}
