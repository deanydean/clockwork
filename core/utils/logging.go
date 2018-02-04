package utils

import (
    "fmt"
    "time"

    "github.com/deanydean/clockwork/core"
    "github.com/deanydean/clockwork/core/triggers"
)

type WatchLogger struct {
    level   int
    handler core.WatchTrigger
}

func GetLogger() *WatchLogger {
    return globalLogger
}

func NewLogger(level int, trigger core.WatchTrigger) *WatchLogger {
    var logger = new(WatchLogger)
    logger.level = level
    logger.handler = trigger
    return logger
}

var LogError = 1 << 0
var LogWarn = 1 << 1
var LogInfo = 1 << 2
var LogDebug = 1 << 4

func (logger *WatchLogger) Error(format string, params ...interface{}) {
    if logger.level >= LogError {
        logger.log(LogError, format, params...)
    }
}

func (logger *WatchLogger) Warn(format string, params ...interface{}) {
    if logger.level >= LogWarn {
        logger.log(LogWarn, format, params...)
    }
}

func (logger *WatchLogger) Info(format string, params ...interface{}) {
    if logger.level >= LogInfo {
        logger.log(LogInfo, format, params...)
    }
}

func (logger *WatchLogger) Debug(format string, params ...interface{}) {
    if logger.level >= LogDebug {
        logger.log(LogDebug, format, params...)
    }
}

var logLevel = "log.level"
var logFormat = "log.format"
var logParams = "log.params"
var logTimestamp = "log.timestamp"

func (logger *WatchLogger) log(level int, format string, params ...interface{}) {
    var logEvent = make(map[string]interface{})
    logEvent[logLevel] = level
    logEvent[logFormat] = format + "\n"
    logEvent[logParams] = params
    logEvent[logTimestamp] = time.Now()

    // Handle log event
    logger.handler.OnEvent(core.NewWatchEvent(logEvent))
}

// Default global log level
var defaultLevel = LogInfo

// Default global log handler
var globalHandler = triggers.NewFuncTrigger(func(event *core.WatchEvent) {
    fmt.Printf(event.GetAsString(logFormat), event.GetAsArray(logParams)...)
})

var globalLogger = NewLogger(defaultLevel, globalHandler)

func SetGlobalLogLevel(level int) {
    globalLogger.level = level
}

func SetGlobalLogHandler(handler core.WatchTrigger) {
    globalLogger.handler = handler
}
