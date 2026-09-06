package main

import (
	"context"
	"errors"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/handlers"
	"microservices-api/internal/loggers"
	"microservices-api/internal/registries"
	"microservices-api/internal/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		appLogger.Error("NewConversionClient returned error", "error", err)
		os.Exit(1)
	}
	defer conversionClient.Close()

	currencyClient, err := clients.NewCurrencyClient(&appConfig)
	if err != nil {
		appLogger.Error("NewCurrencyClient returned error", "error", err)
		os.Exit(1)
	}
	defer currencyClient.Close()

	/* --- --- --- */

	prometheusRegistry := registries.NewPrometheusRegistry()

	/* --- --- --- */

	appServer := &http.Server{
		Addr: ":" + appConfig.App.Port,

		ReadTimeout:  appConfig.App.ReadTimeout,
		IdleTimeout:  appConfig.App.IdleTimeout,
		WriteTimeout: appConfig.App.WriteTimeout,

		Handler: routers.NewAppRouter(&routers.AppRouterHandlers{
			App:               handlers.NewAppHandler(&appConfig, appLogger, currencyClient, conversionClient, prometheusRegistry),
			CurrencyHandler:   handlers.NewCurrencyHandler(&appConfig, appLogger, currencyClient),
			ConversionHandler: handlers.NewConversionHandler(&appConfig, appLogger, conversionClient),
		}, &appConfig, appLogger),
	}

	/* --- --- --- */

	go func() {
		if appConfig.Security.Certificate != "" && appConfig.Security.Key != "" {
			if err := appServer.ListenAndServeTLS(appConfig.Security.Certificate, appConfig.Security.Key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				appLogger.Error("appServer returned error", "error", err)
				os.Exit(1)
			}
		} else {
			if err := appServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				appLogger.Error("appServer returned error", "error", err)
				os.Exit(1)
			}
		}
	}()

	/* --- --- --- */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), appConfig.App.ShutdownTimeout)
	defer cancel()

	if err := appServer.Shutdown(ctx); err != nil {
		appLogger.Error("appServer shutdown failed", "error", err)
		os.Exit(1)
	}
}
