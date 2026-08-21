package configs

import "time"

type AppConfig struct {
	App      App
	Security Security

	HistoryService    HistoryService
	CurrencyService   CurrencyService
	ConversionService ConversionService
}

/* --- --- --- */

type App struct {
	Name string

	Prod bool
	Port string

	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type Security struct {
	Certificate string
	Key         string

	ApiKeys       map[string]bool
	OpenEndpoints map[string]bool
}

/* --- --- --- */

type HistoryService struct {
	Address       string
	Timeout       time.Duration
	Limits        HistoryLimits
	HealthTimeout time.Duration
}

type HistoryLimits struct {
}

/* --- --- --- */

type CurrencyService struct {
	Address       string
	Timeout       time.Duration
	Limits        CurrencyLimits
	HealthTimeout time.Duration
}

type CurrencyLimits struct {
	Rate  Rate
	Rates Rates
}

type Rate struct {
	Limit int
	Burst int
}

type Rates struct {
	Limit int
	Burst int
}

/* --- --- --- */

type ConversionService struct {
	Address       string
	Timeout       time.Duration
	Limits        ConversionLimits
	HealthTimeout time.Duration
}

type ConversionLimits struct {
	Convert Convert
}

type Convert struct {
	Limit int
	Burst int
}
