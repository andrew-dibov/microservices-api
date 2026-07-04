package tools

import (
	"os"
	"strconv"
	"time"
)

func GetStrEnv(key string, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}

func GetDurEnv(key string, def time.Duration) time.Duration {
	if str := os.Getenv(key); str != "" {
		if val, err := time.ParseDuration(str); err == nil {
			return val
		}
	}

	return def
}

func GetIntEnv(key string, def int) int {
	if str := os.Getenv(key); str != "" {
		if val, err := strconv.Atoi(str); err == nil {
			return val
		}
	}
	return def
}
