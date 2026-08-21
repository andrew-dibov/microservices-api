package configs

import (
	"microservices-api/internal/tools"
	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		App: App{
			Port: tools.GetStringEnv("PORT", "8080"),
			Prod: tools.GetBooleanEnv("PROD", false),

			Cert: tools.GetStringEnv("TLS_CERT", ""),
			Key:  tools.GetStringEnv("TLS_KEY", ""),

			Keys: tools.GetStringSetEnv("KEYS", map[string]bool{}),
			Open: tools.GetStringSetEnv("OPEN", map[string]bool{
				"/livez":   true,
				"/readyz":  true,
				"/healthz": true,
				"/metrics": true,
			}),
		},

		Addr: Addr{
			Hist: tools.GetStringEnv("HIST_ADDR", "localhost:50051"),
			Curr: tools.GetStringEnv("CURR_ADDR", "localhost:50052"),
			Conv: tools.GetStringEnv("CONV_ADDR", "localhost:50053"),
		},

		Tout: Tout{
			Hist: tools.GetDurationEnv("HIST_TOUT", 5*time.Second),
			Curr: tools.GetDurationEnv("CURR_TOUT", 5*time.Second),
			Conv: tools.GetDurationEnv("CONV_TOUT", 5*time.Second),

			Read:     tools.GetDurationEnv("READ_TOUT", 10*time.Second),
			Idle:     tools.GetDurationEnv("IDLE_TOUT", 15*time.Second),
			Write:    tools.GetDurationEnv("WRITE_TOUT", 20*time.Second),
			Shutdown: tools.GetDurationEnv("SHUTDOWN_TOUT", 25*time.Second),
		},

		Limt: Limt{
			RateLim: tools.GetIntegerEnv("RATE_LIM", 5),
			RateBur: tools.GetIntegerEnv("RATE_BUR", 10),

			RatesLim: tools.GetIntegerEnv("RATES_LIM", 5),
			RatesBur: tools.GetIntegerEnv("RATES_BUR", 10),

			ConvertLim: tools.GetIntegerEnv("CONVERT_LIM", 5),
			ConvertBur: tools.GetIntegerEnv("CONVERT_BUR", 10),
		},
	}
}
