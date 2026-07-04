package utils

import "time"

func NowStr() string {
	return time.Now().Format("15:04:05")
}

func NowISO() string {
	return time.Now().Format(time.RFC3339)
}
