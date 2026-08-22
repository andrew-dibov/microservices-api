package main

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
	"microservices-api/internal/registries"
	"os"
)

func main() {
	appConfig := configs.NewAppConfig()
	appLogger := loggers.NewAppLogger(appConfig)

	appLogger.Info("config",
		"port", appConfig.App.Port,
		"prod", appConfig.App.Prod,
		"history_address", appConfig.HistoryService.Address,
		"currency_address", appConfig.CurrencyService.Address,
		"conversion_address", appConfig.ConversionService.Address,
	)

	/* --- --- --- */

	conversionClient, err := clients.NewConversionClient(&appConfig)
	if err != nil {
		appLogger.Error("conversionClient returned error", "error", err)
		os.Exit(1)
	}
	defer conversionClient.Close()

	currencyClient, err := clients.NewCurrencyClient(&appConfig)
	if err != nil {
		appLogger.Error("currencyClient returned error", "error", err)
		os.Exit(1)
	}
	defer currencyClient.Close()

	/* --- --- --- */

	prometheusRegistry := registries.NewPrometheusRegistry()

	/* --- --- --- */

	/* --- --- --- */

	// server := &http.Server{
	// 	Addr:         ":" + appConfig.App.Port,
	// 	ReadTimeout:  appConfig.App.ReadTimeout,
	// 	IdleTimeout:  appConfig.App.IdleTimeout,
	// 	WriteTimeout: appConfig.App.WriteTimeout,

	// 	// Handler:,
	// }
}
