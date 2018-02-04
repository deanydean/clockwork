package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deanydean/clockwork/core"
	"github.com/deanydean/clockwork/core/triggers"
	"github.com/deanydean/clockwork/core/utils"
	"github.com/deanydean/clockwork/core/watchers"
	"github.com/deanydean/clockwork/core/watches"
)

var log = utils.GetLogger()

func main() {
	// Get cli flags
	pidFlag := flag.Int("pid", -1, "The pid ")
	debugFlag := flag.Bool("debug", false, "Is debug enabled?")
	flag.Parse()

	if *debugFlag {
		utils.SetGlobalLogLevel(utils.LogDebug)
	}

	var pid = *pidFlag
	if pid <= 0 {
		log.Error("Invalid pid %d\n", pid)
		os.Exit(1)
	} else if !utils.ProcessExists(pid) {
		log.Error("No such process with pid %d\n", pid)
		os.Exit(1)
	}

	var ioWatch = watches.NewProcessHighIOWatch(pid, 10)
	var trigger = triggers.NewFuncTrigger(func(e *core.WatchEvent) {
		fmt.Println("Process", pid, "has high io:")
	})
	var watchMan = watchers.NewWatchMan([]core.Watch{ioWatch})
	watchMan.Watch(trigger)
	select {}
}
