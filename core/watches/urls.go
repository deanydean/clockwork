package watches

import (
	"net/http"
	"time"

	"github.com/deanydean/clockwork/core"
)

var lastModifiedHeader = "Last-Modified"

// URLModifiedTime is a key in WatchEvent for when a URL was last modified
var URLModifiedTime = "url.modifiedtime"

// URLModifiedWatch is a Watch that observes when a URL is modified
type URLModifiedWatch struct {
	url          string
	lastModified time.Time
}

// NewURLModifiedWatch creates a new URLModifiedWatch for the provided url
func NewURLModifiedWatch(url string) *URLModifiedWatch {
	watch := new(URLModifiedWatch)
	watch.url = url

	watch.Observe()

	return watch
}

// Observe whether a URL has changed since the last time it was observed
func (watch *URLModifiedWatch) Observe() *core.WatchEvent {
	resp, err := http.Head(watch.url)

	if err != nil {
		log.Debug("Failed to check url", watch.url, "error", err)
		return nil
	}

	var modTimeStr = resp.Header.Get(lastModifiedHeader)

	if len(modTimeStr) == 0 {
		log.Debug("No modified time for %s cannot observe", watch.url)
		return nil
	}

	var modTime, parseErr = http.ParseTime(modTimeStr)
	if parseErr != nil {
		log.Debug("Failed to parse mod time", modTimeStr, "error", parseErr)
		return nil
	}

	if modTime != watch.lastModified {
		watch.lastModified = modTime
		return core.NewWatchEvent(map[string]interface{}{
			URLModifiedTime: modTime,
		})
	}

	return nil
}
