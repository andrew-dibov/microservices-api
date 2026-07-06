package configs

import "time"

type AppConfig struct {
	Port string

	Prod bool
	Keys map[string]bool
	Open map[string]bool

	Cert string
	Key  string

	Services Services
	Timeouts Timeouts
	Limits   Limits
}

type Services struct {
	Hist string
	Curr string
	Conv string
}

type Timeouts struct {
	Hist time.Duration
	Curr time.Duration
	Conv time.Duration

	Read     time.Duration
	Idle     time.Duration
	Write    time.Duration
	Shutdown time.Duration
}

type Limits struct {
	RateLimit    int
	RatesLimit   int
	ConvertLimit int

	RateBurst    int
	RatesBurst   int
	ConvertBurst int
}
