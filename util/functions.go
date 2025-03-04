package util

import (
	"time"
)

// Uptime */
// Get application uptime in seconds
func Uptime() uint {
	return uint(time.Since(StartTime).Round(time.Second).Seconds())
}
