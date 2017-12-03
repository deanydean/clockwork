package triggers

import "github.com/oddcyborg/watchit/core"

// FuncTrigger
type FuncTrigger struct {
	onEvent func(*core.WatchEvent)
}

// OnEvent
func (trigger FuncTrigger) OnEvent(event *core.WatchEvent) {
	trigger.onEvent(event)
}

// NewFuncTrigger
func NewFuncTrigger(onEvent func(*core.WatchEvent)) *FuncTrigger {
	trigger := new(FuncTrigger)
	trigger.onEvent = onEvent
	return trigger
}

type BroadcastTrigger struct {
	triggers []core.WatchTrigger
}

func (bt BroadcastTrigger) OnEvent(e *core.WatchEvent) {
	for t := range bt.triggers {
		bt.triggers[t].OnEvent(e)
	}
}

func NewBroadcastTrigger(triggers []core.WatchTrigger) *BroadcastTrigger {
	bt := new(BroadcastTrigger)
	bt.triggers = triggers
	return bt
}
