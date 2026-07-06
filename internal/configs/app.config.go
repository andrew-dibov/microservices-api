package configs

import (
	"microservices-api/internal/tools"

	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		Port: tools.GetStrEnv("PORT", "8080"),

		Prod: tools.GetBoolEnv("PROD", false),
		Keys: tools.GetKeysEnv("KEYS", map[string]bool{}),
		Open: tools.GetKeysEnv("OPEN", map[string]bool{
			"/livez":   true,
			"/readyz":  true,
			"/healthz": true,
			"/metrics": true,
		}),

		Cert: tools.GetStrEnv("TLS_CERT", ""),
		Key:  tools.GetStrEnv("TLS_KEY", ""),

		Services: Services{
			Hist: tools.GetStrEnv("HIST_ADDR", "localhost:50051"),
			Curr: tools.GetStrEnv("CURR_ADDR", "localhost:50052"),
			Conv: tools.GetStrEnv("CONV_ADDR", "localhost:50053"),
		},

		Timeouts: Timeouts{
			Hist: tools.GetDurEnv("HIST_TOUT", 5*time.Second),
			Curr: tools.GetDurEnv("CURR_TOUT", 5*time.Second),
			Conv: tools.GetDurEnv("CONV_TOUT", 5*time.Second),

			Read:     tools.GetDurEnv("READ_TOUT", 10*time.Second),
			Idle:     tools.GetDurEnv("IDLE_TOUT", 15*time.Second),
			Write:    tools.GetDurEnv("WRITE_TOUT", 20*time.Second),
			Shutdown: tools.GetDurEnv("SHUTDOWN_TOUT", 25*time.Second),
		},

		Limits: Limits{
			RateLimit:    tools.GetIntEnv("RATE_LIMIT", 5),
			RatesLimit:   tools.GetIntEnv("RATES_LIMIT", 5),
			ConvertLimit: tools.GetIntEnv("CONVERT_LIMIT", 5),

			RateBurst:    tools.GetIntEnv("RATE_BURST", 10),
			RatesBurst:   tools.GetIntEnv("RATES_BURST", 10),
			ConvertBurst: tools.GetIntEnv("CONVERT_BURST", 10),
		},
	}
}
