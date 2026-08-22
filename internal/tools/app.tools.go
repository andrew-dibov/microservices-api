package tools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetStringEnv(key string, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func GetIntegerEnv(key string, def int) int {
	if str := os.Getenv(key); str != "" {
		if val, err := strconv.Atoi(str); err == nil {
			return val
		}
	}
	return def
}

func GetBooleanEnv(key string, def bool) bool {
	if str := os.Getenv(key); str != "" {
		if val, err := strconv.ParseBool(str); err == nil {
			return val
		}
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

func GetStringSetEnv(key string, def map[string]bool) map[string]bool {
	str := os.Getenv(key)
	keys := make(map[string]bool)

	if str != "" {
		for _, k := range strings.Split(str, ",") {
			keys[strings.TrimSpace(k)] = true
		}
	} else {
		keys = def
	}

	return keys
}

/* --- --- --- */

func GenID() string {
	id := make([]byte, 16)

	if _, err := rand.Read(id); err != nil {
		return "00000000000000000000000000000000"
	}

	return hex.EncodeToString(id)
}

func GetID(ctx context.Context) string {
	if val := ctx.Value(ID{}); val != nil {
		return val.(string)
	}

	return ""
}
