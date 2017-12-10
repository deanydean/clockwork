package watches

import (
	"fmt"
	"os"
	"time"

	"github.com/oddcyborg/watchit/core"
)

var FileName = "file.name"

var FileModTime = "file.modifiedtime"

type FileModifiedWatch struct {
	fileName         string
	lastModifiedTime time.Time
}

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
