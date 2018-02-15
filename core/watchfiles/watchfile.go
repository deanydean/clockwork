package watchfiles

import (
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/deanydean/clockwork/core"
	"github.com/deanydean/clockwork/core/triggers"
	"github.com/deanydean/clockwork/core/utils"
	"github.com/deanydean/clockwork/core/watchers"
	"github.com/deanydean/clockwork/core/watches"
)

var log = utils.GetLogger()

// Watchfile containing watch information
type Watchfile struct {
	Watches    []core.Watch
	Triggers   []core.WatchTrigger
	Properties map[string]string
}

// GetWatcherFor gets a watcher for all the watches provided in watchfile
func GetWatcherFor(watchFileName *string) core.Watcher {
	// Get Watchfile
	var watchFile = Load(watchFileName)

	if watchFile == nil {
		log.Warn("Unable to load watchfile=%s", *watchFileName)
		return nil
	}

	// Create a watcher for the watchfile
	return watchers.NewWatchMan(watchFile.Watches)
}

// Load a Watchfile
func Load(watchfile *string) *Watchfile {
	// Read the file
	contents, err := ioutil.ReadFile(*watchfile)
	if err != nil {
		log.Warn("Failed to get watchfile=%s", *watchfile)
		return nil
	}

	// Create a Watchfile object
	var wf = new(Watchfile)

	// Get each lines
	var lines = strings.Split(string(contents), "\n")

	for lineIdx := range lines {
		var line = strings.TrimSpace(lines[lineIdx])

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// Ignore blanks lines and comments
		} else if strings.HasPrefix(line, "WATCH") {
			// Watch defined
			log.Debug("Adding watch from lineNo=%d, line=%s", lineIdx, line)
			wf.Watches = append(wf.Watches, getWatch(line))
		} else if strings.HasPrefix(line, "TELL") {
			// Trigger defined
			log.Debug("Adding trigger from lineNo=%d, line=%s", lineIdx, line)
			wf.Triggers = append(wf.Triggers, getTrigger(line))
		} else if strings.HasPrefix(line, "PROPERTY") {
			// Proprty defined
			log.Debug("Adding property from lineNo=%d, line=%s", lineIdx, line)
			// TODO
		}
	}

	return wf
}

func getWatch(watchLine string) core.Watch {
	// Split line by whitespace
	var sections = strings.Split(watchLine, " ")

	if sections[0] != "WATCH" {
		log.Warn("Unable to parse watchline=%s", watchLine)
		return nil
	}

	// Create the watch from the url
	var url, err = url.Parse(sections[1])
	if err != nil {
		log.Warn("Invalid watch url=%s", sections[1])
		return nil
	}

	switch url.Scheme {
	case "file":
		{
			return watches.NewFileModifiedWatch(url.Path)
		}
	default:
		{
			return watches.NewURLModifiedWatch(sections[1])
		}
	}

	log.Warn("Unknown watch %s", watchLine)
	return nil
}

func getTrigger(tellLine string) core.WatchTrigger {
	// Split line by whitespace
	var sections = strings.Split(tellLine, " ")

	if sections[0] != "TELL" {
		log.Warn("Unable to parse tellline=%s", tellLine)
		return nil
	}

	switch sections[1] {
	case "stdout":
		fallthrough
	default:
		{
			return triggers.NewTextReporterTrigger("")
		}
	}

	log.Warn("Unknown trigger %s", sections[1])
	return nil
}
