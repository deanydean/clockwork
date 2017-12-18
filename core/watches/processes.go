package watches

import (
	"strings"
	"time"

	"github.com/oddcyborg/watchit/core"
	"github.com/oddcyborg/watchit/core/utils"
)

var log = utils.GetLogger()

// ProcessDeathWatch watch that checks if a process has died
type ProcessDeathWatch struct {
	pid int
}

// Observe whether a process has died, returns a WatchEvent if it has, or nil
// if not
func (watch *ProcessDeathWatch) Observe() *core.WatchEvent {
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
func (watch *ProcessHighCPUWatch) Observe() *core.WatchEvent {
	var statsEvent = watch.statsWatch.Observe()

	// Get all the params we need to work out CPU usage
	var uptime = utils.GetSystemUptime()
	var utime = statsEvent.GetAsInteger(StatsProcTime)
	var stime = statsEvent.GetAsInteger(StatsKernTime)
	var cutime = statsEvent.GetAsInteger(StatsProcWaitTime)
	var cstime = statsEvent.GetAsInteger(StatsKernWaitTime)

	var totalTime = utime + stime + cutime + cstime
	var seconds = uptime - float64(watch.procTimeSinceStart)

	var cpuUsage = ((float64(totalTime) / float64(watch.sysClockTick)) /
		seconds)

	statsEvent.Data[StatsCPU] = cpuUsage

	log.Debug("uptime=%d utime=%d stime=%d", uptime, utime, stime)
	log.Debug("total=%d tick=%d seconds=%d usage=%d threshold=%d",
		totalTime, watch.sysClockTick, seconds, cpuUsage, watch.cpuThreshold)

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
	watch.procStartTime = procStats.GetAsInteger(StatsProcStartTime)
	watch.sysClockTick = utils.GetSystemClockTick()

	if watch.sysClockTick == -1 {
		log.Warn("Failed to get system clock speed, cannot create watch")
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
func (watch *ProcessHighMemWatch) Observe() *core.WatchEvent {
	var statsEvent = watch.statsWatch.Observe()

	var rss = statsEvent.GetAsInteger(StatsProcRSS)

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

type ProcessHighIOWatch struct {
	ioThreshold     float64
	lastObservation time.Time
	bytesRead       int64
	bytesWritten    int64
	ioWatch         *ProcessIOWatch
}

func (watch *ProcessHighIOWatch) Observe() *core.WatchEvent {
	var ioEvent = watch.ioWatch.Observe()

	// Work out how much IO the process has performed sine the last check
	var read = ioEvent.GetAsInteger(IOReadBytes)
	var written = ioEvent.GetAsInteger(IOWriteBytes)

	// TODO Check errors
	if read == -1 || written == -1 {
		log.Warn("Failed to get io info")
		return nil
	}

	var now = time.Now()
	var since = now.Sub(watch.lastObservation) / time.Second

	if since > 0 {
		var writesPerSec = (int64(written) - watch.bytesWritten) / int64(since)
		var readsPerSec = (int64(read) - watch.bytesRead) / int64(since)

		log.Debug("rchar=%d sysrc=%d read=%d read/s=%d",
			ioEvent.GetAsInteger(IOReadChar),
			ioEvent.GetAsInteger(IOReadCalls),
			read, readsPerSec)
		log.Debug("wchar=%d syswc=%d written=%d read/s=%d",
			ioEvent.GetAsInteger(IOWriteChar),
			ioEvent.GetAsInteger(IOWriteCalls),
			written, writesPerSec)

		ioEvent.Data[IOReadsPerSec] = readsPerSec
		ioEvent.Data[IOWritesPerSec] = writesPerSec
	}

	watch.lastObservation = now
	watch.bytesRead = int64(read)
	watch.bytesWritten = int64(written)

	// Nothing to report
	return nil
}

func NewProcessHighIOWatch(pid int, threshold float64) *ProcessHighIOWatch {
	watch := new(ProcessHighIOWatch)
	watch.ioThreshold = threshold
	watch.ioWatch = new(ProcessIOWatch)
	watch.ioWatch.pid = pid

	// Init the watch with an initial value
	watch.bytesRead = 0
	watch.bytesWritten = 0
	watch.lastObservation = time.Now()
	watch.Observe()

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
func (watch *ProcessStatsWatch) Observe() *core.WatchEvent {
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

type ProcessIOWatch struct {
	pid int
}

var IORaw = "io.raw"
var IOReadChar = "io.rchar"
var IOWriteChar = "io.wchar"
var IOReadCalls = "io.syscr"
var IOWriteCalls = "io.syscw"
var IOReadBytes = "io.read_bytes"
var IOWriteBytes = "io.write_bytes"
var IOCancelledWriteBytes = "io.cancelled_write_bytes"
var IOWritesPerSec = "io.writes_per_sec"
var IOReadsPerSec = "io.reads_per_sec"

func (watch *ProcessIOWatch) Observe() *core.WatchEvent {
	var ioStr, err = utils.GetProcessIO(watch.pid)

	if err != nil {
		return nil
	}

	var io = make(map[string]interface{})
	var rawIO = strings.Split(ioStr, "\n")
	io[IORaw] = rawIO

	// Parse the lines, using the key for the event key
	for _, line := range rawIO {
		var kv = strings.Split(line, ":")
		if len(kv) == 2 {
			io["io."+kv[0]] = strings.TrimSpace(kv[1])
		}
	}

	return core.NewWatchEvent(io)
}
