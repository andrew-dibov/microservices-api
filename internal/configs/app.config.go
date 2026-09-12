package configs

import (
	"microservices-api/internal/modules"
	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		App: App{
			Name: modules.GetStringEnv("APP_NAME", "microservices-api"),

			Prod: modules.GetBooleanEnv("APP_PROD", false),
			Port: modules.GetStringEnv("APP_PORT", "8080"),

			ReadyzTimeout: modules.GetDurationEnv("APP_READYZ_TIMEOUT", 2*time.Second),

			ReadTimeout:  modules.GetDurationEnv("APP_READ_TIMEOUT", 5*time.Second),
			IdleTimeout:  modules.GetDurationEnv("APP_IDLE_TIMEOUT", 5*time.Second),
			WriteTimeout: modules.GetDurationEnv("APP_WRITE_TIMEOUT", 5*time.Second),

			ShutdownTimeout: modules.GetDurationEnv("APP_SHUTDOWN_TIMEOUT", 5*time.Second),
		},

		Security: Security{
			Certificate: modules.GetStringEnv("SECURITY_CERTIFICATE", ""),
			Key:         modules.GetStringEnv("SECURITY_KEY", ""),

			ApiKeys: modules.GetStringSetEnv("SECURITY_API_KEYS", map[string]bool{}),
			OpenEndpoints: modules.GetStringSetEnv("SECURITY_OPEN_ENDPOINTS", map[string]bool{
				"/livez":   true,
				"/readyz":  true,
				"/healthz": true,
				"/metrics": true,
			}),
		},

		HistoryService: HistoryService{
			Address:       modules.GetStringEnv("HISTORY_ADDRESS", "localhost:50051"),
			Timeout:       modules.GetDurationEnv("HISTORY_TIMEOUT", 5*time.Second),
			HealthTimeout: modules.GetDurationEnv("HISTORY_HEALTH_TIMEOUT", 2*time.Second),

			Limits: HistoryLimits{},
		},

		CurrencyService: CurrencyService{
			Address:       modules.GetStringEnv("CURRENCY_ADDRESS", "localhost:50052"),
			Timeout:       modules.GetDurationEnv("CURRENCY_TIMEOUT", 5*time.Second),
			HealthTimeout: modules.GetDurationEnv("CURRENCY_HEALTH_TIMEOUT", 2*time.Second),

			Limits: CurrencyLimits{
				Rate: Rate{
					Limit: modules.GetIntegerEnv("CURRENCY_RATE_LIMIT", 5),
					Burst: modules.GetIntegerEnv("CURRENCY_RATE_BURST", 10),
				},
				Rates: Rates{
					Limit: modules.GetIntegerEnv("CURRENCY_RATES_LIMIT", 5),
					Burst: modules.GetIntegerEnv("CURRENCY_RATES_BURST", 10),
				},
			},
		},

		ConversionService: ConversionService{
			Address:       modules.GetStringEnv("CONVERSION_ADDRESS", "localhost:50053"),
			Timeout:       modules.GetDurationEnv("CONVERSION_TIMEOUT", 5*time.Second),
			HealthTimeout: modules.GetDurationEnv("CONVERSION_HEALTH_TIMEOUT", 2*time.Second),

			Limits: ConversionLimits{
				Convert: Convert{
					Limit: modules.GetIntegerEnv("CONVERSION_CONVERT_LIMIT", 5),
					Burst: modules.GetIntegerEnv("CONVERSION_CONVERT_BURST", 10),
				},
			},
		},
	}
}
