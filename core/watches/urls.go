package watches

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oddcyborg/watchit/core"
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

	event := watch.Observe()
	if event == nil {
		fmt.Println("Failed to observe URL", url)
		return nil
	}

	return watch
}

// Observe whether a URL has changed since the last time it was observed
func (watch *URLModifiedWatch) Observe() *core.WatchEvent {
	resp, err := http.Head(watch.url)

	if err != nil {
		fmt.Println("Failed to check url", watch.url, "error", err)
		return nil
	}

	var modTimeStr = resp.Header.Get(lastModifiedHeader)

	if len(modTimeStr) == 0 {
		fmt.Println("No modified time for", watch.url, "cannot observe")
		return nil
	}

	var modTime, parseErr = http.ParseTime(modTimeStr)
	if parseErr != nil {
		fmt.Println("Failed to parse mod time", modTimeStr, "error", parseErr)
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
