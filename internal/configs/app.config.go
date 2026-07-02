package configs

import (
	"microservices-api/internal/tools"

	"time"
)

type AppConfig struct {
	Port string

	Services Services
	Timeouts Timeouts
}

type Services struct {
	History    string
	Currency   string
	Conversion string
}

type Timeouts struct {
	Read     time.Duration
	Idle     time.Duration
	Write    time.Duration
	Shutdown time.Duration
}

func NewAppConfig() AppConfig {
	return AppConfig{
		Port: tools.GetEnv("PORT", "8080"),

		Services: Services{
			History:    tools.GetEnv("HISTORY_SERVICE", "localhost:50051"),
			Currency:   tools.GetEnv("CURRENCY_SERVICE", "localhost:50052"),
			Conversion: tools.GetEnv("CONVERSION_SERVICE", "localhost:50053"),
		},

		Timeouts: Timeouts{
			Read:     tools.GetDurationEnv("READ_TIMEOUT", 10*time.Second),
			Idle:     tools.GetDurationEnv("IDLE_TIMEOUT", 15*time.Second),
			Write:    tools.GetDurationEnv("WRITE_TIMEOUT", 20*time.Second),
			Shutdown: tools.GetDurationEnv("SHUTDOWN_TIMEOUT", 25*time.Second),
		},
	}
}
