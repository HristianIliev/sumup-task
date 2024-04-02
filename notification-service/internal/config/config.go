package config

import (
	"os"
)

var Config = func(key string) string {
	switch key {
	case "API_PORT":
		return os.Getenv("API_PORT")
	case "TABLE_NAME":
		return os.Getenv("TABLE_NAME")
	case "SNS_TOPIC":
		return os.Getenv("SNS_TOPIC")
	default:
		return ""
	}
}
