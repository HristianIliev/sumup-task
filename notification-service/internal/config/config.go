// nolint:depguard
package config

import (
	"os"
)

var Config = func(key string) string {
	switch key {
	case "API_PORT":
		return os.Getenv("API_PORT")
	case "DB_CONNECTION_URL":
		return os.Getenv("DB_CONNECTION_URL")
	default:
		return ""
	}
}
