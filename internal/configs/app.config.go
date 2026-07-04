package configs

import (
	"microservices-api/internal/tools"

	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		Port: tools.GetStrEnv("PORT", "8080"),

		Services: Services{
			History:    tools.GetStrEnv("HISTORY_SERVICE", "localhost:50051"),
			Currency:   tools.GetStrEnv("CURRENCY_SERVICE", "localhost:50052"),
			Conversion: tools.GetStrEnv("CONVERSION_SERVICE", "localhost:50053"),
		},

		Timeouts: Timeouts{
			Read:     tools.GetDurEnv("READ_TIMEOUT", 10*time.Second),
			Idle:     tools.GetDurEnv("IDLE_TIMEOUT", 15*time.Second),
			Write:    tools.GetDurEnv("WRITE_TIMEOUT", 20*time.Second),
			Shutdown: tools.GetDurEnv("SHUTDOWN_TIMEOUT", 25*time.Second),

			History:    tools.GetDurEnv("HISTORY_TIMEOUT", 5*time.Second),
			Currency:   tools.GetDurEnv("CURRENCY_TIMEOUT", 5*time.Second),
			Conversion: tools.GetDurEnv("CONVERSION_TIMEOUT", 5*time.Second),
		},

		Limits: Limits{
			RateLimit:    tools.GetIntEnv("RATE_LIMIT", 5),
			RateBurst:    tools.GetIntEnv("RATE_BURST", 10),
			RatesLimit:   tools.GetIntEnv("RATES_LIMIT", 5),
			RatesBurst:   tools.GetIntEnv("RATES_BURST", 10),
			ConvertLimit: tools.GetIntEnv("CONVERT_LIMIT", 5),
			ConvertBurst: tools.GetIntEnv("CONVERT_BURST", 10),
		},
	}
}
