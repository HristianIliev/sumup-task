package config

import (
	"os"
)

var Config = func(key string) string {
	switch key {
	case "QUEUE_URL":
		return os.Getenv("QUEUE_URL")
	case "WORKERS":
		return os.Getenv("WORKERS")
	case "POLL_DELAY":
		return os.Getenv("POLL_DELAY")
	default:
		return ""
	}
}
