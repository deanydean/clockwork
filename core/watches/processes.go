package watches

import (
	"strconv"
	"strings"
	"time"

	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/utils"
)

type ProcessDeathWatch struct {
	pid int
}

func (watch ProcessDeathWatch) Observe() *core.WatchEvent {
	if !utils.ProcessExists(watch.pid) {
		return core.NewWatchEvent(nil)
	}

	return nil
}

func NewProcessDeathWatch(pid int) *ProcessDeathWatch {
	watch := new(ProcessDeathWatch)
	watch.pid = pid
	return watch
}

// ProcessHighCPUWatch will watch for a process CPU going over a threshold
type ProcessHighCPUWatch struct {
	pid          int
	cpuThreshold float64
	statsWatch   ProcessStatsWatch
}

func (watch ProcessHighCPUWatch) Observe() *core.WatchEvent {
	return nil
}

func NewProcessHighCPUWatch(pid int, threshold float64) *ProcessHighCPUWatch {
	watch := new(ProcessHighCPUWatch)
	watch.pid = pid
	watch.cpuThreshold = threshold
	return watch
}

// ProcessHighMemWatch will watch for a process memory going over a thresold
type ProcessHighMemWatch struct {
	pid          int
	memThreshold float64
	statsWatch   ProcessStatsWatch
}

func (watch ProcessHighMemWatch) Observe() *core.WatchEvent {
	return nil
}

func NewProcessHighMemWatch(pid int, threshold float64) *ProcessHighMemWatch {
	watch := new(ProcessHighMemWatch)
	watch.pid = pid
	watch.memThreshold = threshold
	return watch
}

// TODO switch stats back to stats object/struct

var STATS_RAW = "stats.raw"
var STATS_CPU = "stats.cpu"
var STATS_MEM = "stats.mem"
var STATS_TIMESTAMP = "stats.timestamp"
var STATS_PROCTIME = "stats.proctime"
var STATS_PROCMEM = "stats.procmem"

type ProcessStatsWatch struct {
	pid int
}

// Observe process stats for the watch's pid
func (watch ProcessStatsWatch) Observe() *core.WatchEvent {
	var statFile = "/proc/" + strconv.Itoa(watch.pid) + "/stat"
	fileStr, err := utils.GetFileAsString(statFile)
	var timestamp = time.Now()

	if err != nil {
		return nil
	}

	// Parse stats file and put in the defaults
	var stats = make(map[string]interface{})
	var rawStats = strings.Split(fileStr, " ")
	stats[STATS_RAW] = rawStats
	stats[STATS_CPU] = -1
	stats[STATS_MEM] = -1
	stats[STATS_TIMESTAMP] = timestamp

	// Put in the specifics
	stats[STATS_PROCTIME] = rawStats[14]
	stats[STATS_PROCMEM] = rawStats[24]

	return core.NewWatchEvent(stats)
}
