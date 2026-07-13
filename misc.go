package main

import "time"

// GetTimeMs returns the current wall-clock time in milliseconds, used for
// search time management.
func GetTimeMs() int64 {
	return time.Now().UnixMilli()
}
