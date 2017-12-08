package watches

import (
	"fmt"
	"strings"
	"time"

	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/core/utils"
)

// ProcessDeathWatch watch that checks if a process has died
type ProcessDeathWatch struct {
	pid int
}

// Observe whether a process has died, returns a WatchEvent if it has, or nil
// if not
func (watch ProcessDeathWatch) Observe() *core.WatchEvent {
	if !utils.ProcessExists(watch.pid) {
		return core.NewWatchEvent(nil)
	}

	return nil
}

// NewProcessDeathWatch returns a new ProcessDeathWatch for the provided pid
func NewProcessDeathWatch(pid int) *ProcessDeathWatch {
	watch := new(ProcessDeathWatch)
	watch.pid = pid
	return watch
}

// ProcessHighCPUWatch will watch for a process CPU going over a threshold
type ProcessHighCPUWatch struct {
	cpuThreshold       float64
	procStartTime      int
	sysClockTick       int
	procTimeSinceStart int
	statsWatch         *ProcessStatsWatch
}

// Observe whether a process CPU is high, returns a WatchEvent if it is, or nil
// if it's not
func (watch ProcessHighCPUWatch) Observe() *core.WatchEvent {
	var statsEvent = watch.statsWatch.Observe()

	// Get all the params we need to work out CPU usage
	var uptime = utils.GetSystemUptime()
	var utime, _ = statsEvent.GetAsInteger(StatsProcTime)
	var stime, _ = statsEvent.GetAsInteger(StatsKernTime)
	var cutime, _ = statsEvent.GetAsInteger(StatsProcWaitTime)
	var cstime, _ = statsEvent.GetAsInteger(StatsKernWaitTime)

	var totalTime = utime + stime + cutime + cstime
	var seconds = uptime - float64(watch.procTimeSinceStart)

	var cpuUsage = ((float64(totalTime) / float64(watch.sysClockTick)) /
		seconds)

	statsEvent.Data[StatsCPU] = cpuUsage

	fmt.Println("uptime", uptime, "utime", utime, "stime", stime)
	fmt.Println("total", totalTime, "tick", watch.sysClockTick,
		"seconds", seconds, "usage", cpuUsage, "threshold", watch.cpuThreshold)

	if cpuUsage > watch.cpuThreshold {
		return statsEvent
	}

	return nil
}

// NewProcessHighCPUWatch returns a new ProcessHighCPUWatch for the provided pid
// with the provided high CPU threshold
func NewProcessHighCPUWatch(pid int, threshold float64) *ProcessHighCPUWatch {
	watch := new(ProcessHighCPUWatch)
	watch.cpuThreshold = threshold
	watch.statsWatch = new(ProcessStatsWatch)
	watch.statsWatch.pid = pid

	var procStats = watch.statsWatch.Observe()
	watch.procStartTime, _ = procStats.GetAsInteger(StatsProcStartTime)
	watch.sysClockTick = utils.GetSystemClockTick()

	if watch.sysClockTick == -1 {
		fmt.Println("Failed to get system clock speed, cannot create watch")
		return nil
	}

	watch.procTimeSinceStart = watch.procStartTime / watch.sysClockTick

	return watch
}

// ProcessHighMemWatch will watch for a process memory going over a thresold
type ProcessHighMemWatch struct {
	memThreshold float64
	statsWatch   *ProcessStatsWatch
}

// Observe whether a process has high memory usage, returns a WatchEvent if it
// has or nil if it hasn't
func (watch ProcessHighMemWatch) Observe() *core.WatchEvent {
	var statsEvent = watch.statsWatch.Observe()

	var rss, _ = statsEvent.GetAsInteger(StatsProcRSS)

	var bytesInUse = rss * utils.GetPageSize()
	statsEvent.Data[StatsMem] = bytesInUse

	if float64(bytesInUse) > watch.memThreshold {
		return statsEvent
	}

	// Nothing to report
	return nil
}

// NewProcessHighMemWatch returns a new ProcessHighMemWatch for the provided pid
// with the provided high memory threshold
func NewProcessHighMemWatch(pid int, threshold float64) *ProcessHighMemWatch {
	watch := new(ProcessHighMemWatch)
	watch.memThreshold = threshold
	watch.statsWatch = new(ProcessStatsWatch)
	watch.statsWatch.pid = pid
	return watch
}

// StatsRaw is a key in WatchEvent for raw process information
var StatsRaw = "stats.raw"

// StatsCPU is a key in WatchEvent for process CPU usage
var StatsCPU = "stats.cpu"

// StatsMem is a key in WatchEvent for process memory usage
var StatsMem = "stats.mem"

// StatsTimestamp is a key in WatchEvent for process timestamp
var StatsTimestamp = "stats.timestamp"

// StatsProcTime is a key in WatchEvent for process processor time
var StatsProcTime = "stats.proctime"

// StatsProcRSS is a key in WatchEvent for process rss value
var StatsProcRSS = "stats.procrss"

// StatsKernTime is a key in WatchEvent for process kernel time
var StatsKernTime = "stats.kerntime"

// StatsProcWaitTime is a key in WatchEvent for process user wait time
var StatsProcWaitTime = "stats.procwaittime"

// StatsKernWaitTime is a key in WatchEvent for process kernel wait time
var StatsKernWaitTime = "stats.kernwaittime"

// StatsProcStartTime is a key in WatchEvent for process start time
var StatsProcStartTime = "stats.procstarttime"

// ProcessStatsWatch watch that observes process information
type ProcessStatsWatch struct {
	pid int
}

// Observe process stats for the watch's pid
func (watch ProcessStatsWatch) Observe() *core.WatchEvent {
	var statsStr, err = utils.GetProcessStats(watch.pid)

	if err != nil {
		return nil
	}

	// Parse stats file and put in the defaults
	var stats = make(map[string]interface{})
	var rawStats = strings.Split(statsStr, " ")
	stats[StatsRaw] = rawStats
	stats[StatsTimestamp] = time.Now()

	// Put in the specifics
	stats[StatsProcTime] = rawStats[13]
	stats[StatsKernTime] = rawStats[14]
	stats[StatsProcWaitTime] = rawStats[15]
	stats[StatsKernWaitTime] = rawStats[16]
	stats[StatsProcStartTime] = rawStats[21]
	stats[StatsProcRSS] = rawStats[23]

	return core.NewWatchEvent(stats)
}
