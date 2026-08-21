package configs

import (
	"microservices-api/internal/tools"
	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		App: App{
			Name: tools.GetStringEnv("APP_NAME", "microservices-api"),

			Prod: tools.GetBooleanEnv("APP_PROD", false),
			Port: tools.GetStringEnv("APP_PORT", "8080"),

			ReadTimeout:     tools.GetDurationEnv("APP_READ_TIMEOUT", 5*time.Second),
			IdleTimeout:     tools.GetDurationEnv("APP_IDLE_TIMEOUT", 5*time.Second),
			WriteTimeout:    tools.GetDurationEnv("APP_WRITE_TIMEOUT", 5*time.Second),
			ShutdownTimeout: tools.GetDurationEnv("APP_SHUTDOWN_TIMEOUT", 5*time.Second),
		},

		Security: Security{
			Certificate: tools.GetStringEnv("SECURITY_CERTIFICATE", ""),
			Key:         tools.GetStringEnv("SECURITY_KEY", ""),

			ApiKeys: tools.GetStringSetEnv("SECURITY_API_KEYS", map[string]bool{}),
			OpenEndpoints: tools.GetStringSetEnv("SECURITY_OPEN_ENDPOINTS", map[string]bool{
				"/livez":   true,
				"/readyz":  true,
				"/healthz": true,
				"/metrics": true,
			}),
		},

		HistoryService: HistoryService{
			Address:       tools.GetStringEnv("HISTORY_ADDRESS", "localhost:50051"),
			Timeout:       tools.GetDurationEnv("HISTORY_TIMEOUT", 5*time.Second),
			HealthTimeout: tools.GetDurationEnv("HISTORY_HEALTH_TIMEOUT", 2*time.Second),

			Limits: HistoryLimits{},
		},

		CurrencyService: CurrencyService{
			Address:       tools.GetStringEnv("CURRENCY_ADDRESS", "localhost:50052"),
			Timeout:       tools.GetDurationEnv("CURRENCY_TIMEOUT", 5*time.Second),
			HealthTimeout: tools.GetDurationEnv("CURRENCY_HEALTH_TIMEOUT", 2*time.Second),

			Limits: CurrencyLimits{
				Rate: Rate{
					Limit: tools.GetIntegerEnv("CURRENCY_RATE_LIMIT", 5),
					Burst: tools.GetIntegerEnv("CURRENCY_RATE_BURST", 10),
				},
				Rates: Rates{
					Limit: tools.GetIntegerEnv("CURRENCY_RATES_LIMIT", 5),
					Burst: tools.GetIntegerEnv("CURRENCY_RATES_BURST", 10),
				},
			},
		},

		ConversionService: ConversionService{
			Address:       tools.GetStringEnv("CONVERSION_ADDRESS", "localhost:50053"),
			Timeout:       tools.GetDurationEnv("CONVERSION_TIMEOUT", 5*time.Second),
			HealthTimeout: tools.GetDurationEnv("CONVERSION_HEALTH_TIMEOUT", 2*time.Second),

			Limits: ConversionLimits{
				Convert: Convert{
					Limit: tools.GetIntegerEnv("CONVERSION_CONVERT_LIMIT", 5),
					Burst: tools.GetIntegerEnv("CONVERSION_CONVERT_BURST", 10),
				},
			},
		},
	}
}
