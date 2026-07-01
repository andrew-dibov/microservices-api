package tools

import (
	"os"
	"time"
)

func GetEnv(key string, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}

func GetDurationEnv(key string, def time.Duration) time.Duration {
	if str := os.Getenv(key); str != "" {
		if val, err := time.ParseDuration(str); err == nil {
			return val
		}
	}

	return def
}
