package configs

import "time"

type AppConfig struct {
	Prod bool
	Port string

	Keys map[string]bool

	Services Services
	Timeouts Timeouts
	Limits   Limits
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

	History    time.Duration
	Currency   time.Duration
	Conversion time.Duration
}

type Limits struct {
	RateLimit    int
	RateBurst    int
	RatesLimit   int
	RatesBurst   int
	ConvertLimit int
	ConvertBurst int
}
