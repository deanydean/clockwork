package utils

import "strconv"

// ProcessExists returns true if a process identified by pid exists, false if
// not
func ProcessExists(pid int) bool {
	return PathExists("/proc/" + strconv.Itoa(pid))
}
