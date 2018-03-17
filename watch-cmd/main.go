package main

import (
	"flag"
	"fmt"

	"github.com/deanydean/clockwork/core"
	"github.com/deanydean/clockwork/core/triggers"
	"github.com/deanydean/clockwork/core/utils"
	"github.com/deanydean/clockwork/core/watchers"
	"github.com/deanydean/clockwork/core/watches"
)

var log = utils.GetLogger()

func main() {
	// Get cli flags
	debugFlag := flag.Bool("debug", false, "Is debug enabled?")
	flag.Parse()

	var i = 0
	if *debugFlag {
		i += 1
		utils.SetGlobalLogLevel(utils.LogDebug)
	}

	var cmdLine = flag.Args()
	var cmd = cmdLine[0]
	var args = cmdLine[1:]

	var watch = watches.NewCommandWatch(cmd, args)
	var watchMan = watchers.NewWatchMan([]core.Watch{watch})

	// Create the triggers
	var outputTrigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		log.Debug("Event: code=%d stop?%b", e.Result(), e.ShouldStop())

		if e.Get("output-error") == nil {
			fmt.Println("> ", e.Get("output"))
		} else {
			fmt.Println("Failed:", e.Get("output-error-message"))
			fmt.Println("ERROR:", e.Get("output-error"))
		}
		if e.Get("error-error") == nil {
			fmt.Println("X ", e.Get("error"))
		} else {
			fmt.Println("Failed:", e.Get("error-error-message"))
			fmt.Println("ERROR:", e.Get("error-error"))
		}
	})

	// Start watching
	log.Info("Watching %s", cmd)
	watchMan.Watch(outputTrigger)

	select {}
}
