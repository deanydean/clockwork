package main

import (
    "fmt"
    "os"
    "path"
    "reflect"

    "github.com/deanydean/clockwork/core/utils"
    "github.com/deanydean/clockwork/core/watchfiles"
)

var log = utils.GetLogger()

func main() {
    var cwd, err = os.Getwd()
    if err != nil {
        fmt.Println("Can't to get working directory, cannot read Watchfile")
        return
    }

    // Trace all paths
    utils.SetGlobalLogLevel(utils.LogDebug)

    var watchFileName = path.Join(cwd, "Watchfile")
    var watchFile = watchfiles.Load(&watchFileName)

    if watchFile == nil {
        log.Warn("Unable to load watchfile=%s", watchFileName)
        return
    }

    // Report on status
    for w := range watchFile.Watches {
        log.Info("Watch %s loaded", reflect.TypeOf(watchFile.Watches[w]))
    }
    for t := range watchFile.Triggers {
        log.Info("Trigger %s loaded", reflect.TypeOf(watchFile.Triggers[t]))
    }
    for property := range watchFile.Properties {
        log.Info("Using property %s=%s", property, watchFile.Properties[property])
    }
}
