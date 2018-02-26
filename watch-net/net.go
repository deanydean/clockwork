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
	layerParam := flag.String("layer", "IPv4", "The network layer to watch")
	ifaceParam := flag.String("interface", "eth0", "The network interface to watch")
	debugParam := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	log.Info("Running network watch")

	if *debugParam {
		utils.SetGlobalLogLevel(utils.LogDebug)
		log.Debug("Debugging enabled")
	}

	// Create the watch
	var modifiedWatch = watches.NewNetWatch(layerParam, ifaceParam)
	if modifiedWatch == nil {
		fmt.Println("Failed to create net watch")
		return
	}

	log.Debug("Created watch, now creating watchman")

	var watchMan = watchers.NewWatchMan([]core.Watch{modifiedWatch})

	// Create the triggers
	var modifiedTrigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Printf("Some net event %s", e.Data)
	})

	// Start watching
	log.Info("Starting to watch....")
	watchMan.Watch(modifiedTrigger)

	select {}
}
