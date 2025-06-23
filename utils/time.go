package utils

import (
	"time"
)

// GetNow returns current time in default layout
func GetNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
