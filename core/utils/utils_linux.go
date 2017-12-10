package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ProcessExists returns true if a process identified by pid exists, false if
// not
func ProcessExists(pid int) bool {
	return PathExists("/proc/" + strconv.Itoa(pid))
}

// GetProcessStats returns the process stats information for the provided pid
// and an error which is set if something went wrong
func GetProcessStats(pid int) (string, error) {
	var statFile = "/proc/" + strconv.Itoa(pid) + "/stat"
	return GetFileAsString(statFile)
}

func GetSystemUptime() float64 {
	var statFile = "/proc/uptime"
	var fileContents, _ = GetFileAsString(statFile)
	var uptime, _ = strconv.ParseFloat(strings.Split(fileContents, " ")[0], 64)
	return uptime
}

func GetSystemClockTick() int {
	var ticks, err = exec.Command("getconf", "CLK_TCK").Output()

	if err == nil {
		var tickString = strings.TrimSpace(string(ticks))
		var tickValue, numErr = strconv.Atoi(tickString)

		if numErr != nil {
			fmt.Println("Failed to parse tick string=", tickString,
				" err=", numErr)
			return -1
		}

		return tickValue
	}

	fmt.Println("Failed to get system clock speed err=", err)
	return -1
}

func GetPageSize() int {
	var pageSizeBytes, err = exec.Command("getconf", "PAGE_SIZE").Output()

	if err == nil {
		var pageSizeString = strings.TrimSpace(string(pageSizeBytes))
		var pageSize, numErr = strconv.Atoi(pageSizeString)

		if numErr != nil {
			fmt.Println("Failed to parse page size string=", pageSizeString,
				" err=", numErr)
			return -1
		}

		return pageSize
	}

	fmt.Println("Failed to get system page size err=", err)
	return -1
}
