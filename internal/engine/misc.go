package engine

import "time"

// GetTimeMs returns the current wall-clock time in milliseconds, used for
// search time management.
func GetTimeMs() int64 {
	return time.Now().UnixMilli()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
