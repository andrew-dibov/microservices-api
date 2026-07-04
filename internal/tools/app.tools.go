package tools

import (
	"os"
	"strconv"
	"strings"
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

func GetBolEnv(key string, def bool) bool {
	if str := os.Getenv(key); str != "" {
		if val, err := strconv.ParseBool(str); err == nil {
			return val
		}
	}
	return def
}

func GetKysEnv(key string, def map[string]bool) map[string]bool {
	str := os.Getenv(key)
	kys := make(map[string]bool)

	if str != "" {
		for _, key := range strings.Split(str, ",") {
			kys[strings.TrimSpace(key)] = true
		}
	} else {
		kys = def
	}

	return kys
}
