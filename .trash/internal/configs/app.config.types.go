package configs

import "time"

type AppConfig struct {
	App  App
	Addr Addr
	Tout Tout
	Limt Limt
}

type App struct {
	Port string
	Prod bool

	Cert string
	Key  string

	Keys map[string]bool
	Open map[string]bool
}

type Addr struct {
	Hist string
	Curr string
	Conv string
}

type Tout struct {
	Hist time.Duration
	Curr time.Duration
	Conv time.Duration

	Read     time.Duration
	Idle     time.Duration
	Write    time.Duration
	Shutdown time.Duration
}

type Limt struct {
	RateLim    int
	RatesLim   int
	ConvertLim int

	RateBur    int
	RatesBur   int
	ConvertBur int
}
