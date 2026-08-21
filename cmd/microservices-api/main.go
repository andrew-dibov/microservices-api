package main

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/loggers"
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
}
