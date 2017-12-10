package watches

import (
	"fmt"
	"os"
	"time"

	"github.com/oddcyborg/watchit/core"
)

// FileName is a key in EventWatch for a file name
var FileName = "file.name"

// FileModTime is a key in EventWatch for a file's modified time
var FileModTime = "file.modifiedtime"

// FileModifiedWatch is a Watch that observes when a file modified time changes
type FileModifiedWatch struct {
	fileName         string
	lastModifiedTime time.Time
}

// NewFileModifiedWatch creates a new FileModifiedWatch for the provided file
func NewFileModifiedWatch(file string) *FileModifiedWatch {
	watch := new(FileModifiedWatch)
	watch.fileName = file

	// Get file information
	var info, err = os.Stat(watch.fileName)

	if err != nil {
		fmt.Println("Failed to read file", watch.fileName, "err", err)
		return nil
	}

	watch.lastModifiedTime = info.ModTime()
	return watch
}

// Observe whether a file has been modified since it was last observed
func (watch *FileModifiedWatch) Observe() *core.WatchEvent {
	// Get file information
	var info, err = os.Stat(watch.fileName)

	if err != nil {
		fmt.Println("Failed to read file", watch.fileName, "err", err)
		return nil
	}

	if info.ModTime() != watch.lastModifiedTime {
		watch.lastModifiedTime = info.ModTime()
		return core.NewWatchEvent(map[string]interface{}{
			FileName:    watch.fileName,
			FileModTime: info.ModTime(),
		})
	}

	return nil
}
