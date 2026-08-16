package configs

import (
	"microservices-api/internal/tools"
	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		App: App{
			Port: tools.GetStrEnv("PORT", "8080"),
			Prod: tools.GetBoolEnv("PROD", false),

			Cert: tools.GetStrEnv("TLS_CERT", ""),
			Key:  tools.GetStrEnv("TLS_KEY", ""),

			Keys: tools.GetKeysEnv("KEYS", map[string]bool{}),
			Open: tools.GetKeysEnv("OPEN", map[string]bool{
				"/livez":   true,
				"/readyz":  true,
				"/healthz": true,
				"/metrics": true,
			}),
		},

		Addr: Addr{
			Hist: tools.GetStrEnv("HIST_ADDR", "localhost:50051"),
			Curr: tools.GetStrEnv("CURR_ADDR", "localhost:50052"),
			Conv: tools.GetStrEnv("CONV_ADDR", "localhost:50053"),
		},

		Tout: Tout{
			Hist: tools.GetDurEnv("HIST_TOUT", 5*time.Second),
			Curr: tools.GetDurEnv("CURR_TOUT", 5*time.Second),
			Conv: tools.GetDurEnv("CONV_TOUT", 5*time.Second),

			Read:     tools.GetDurEnv("READ_TOUT", 10*time.Second),
			Idle:     tools.GetDurEnv("IDLE_TOUT", 15*time.Second),
			Write:    tools.GetDurEnv("WRITE_TOUT", 20*time.Second),
			Shutdown: tools.GetDurEnv("SHUTDOWN_TOUT", 25*time.Second),
		},

		Limt: Limt{
			RateLim: tools.GetIntEnv("RATE_LIM", 5),
			RateBur: tools.GetIntEnv("RATE_BUR", 10),

			RatesLim: tools.GetIntEnv("RATES_LIM", 5),
			RatesBur: tools.GetIntEnv("RATES_BUR", 10),

			ConvertLim: tools.GetIntEnv("CONVERT_LIM", 5),
			ConvertBur: tools.GetIntEnv("CONVERT_BUR", 10),
		},
	}
}
