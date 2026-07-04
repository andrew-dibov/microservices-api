package main

import (
	"log/slog"
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/handlers"
	"microservices-api/internal/routers"

	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := configs.NewAppConfig()

	var hand slog.Handler
	if cfg.Prod {
		hand = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		hand = slog.NewTextHandler(os.Stdout, nil)
	}

	log := slog.New(hand)
	log.Info("app config loaded", "port", cfg.Port, "history", cfg.Services.History, "currency", cfg.Services.Currency, "conversion", cfg.Services.Conversion)

	currency, err := clients.NewCurrencyClient(cfg.Services.Currency, cfg.Timeouts.Currency)
	if err != nil {
		log.Error("currency client failed", "error", err)
		os.Exit(1)
	}
	defer currency.Close()

	conversion, err := clients.NewConversionClient(cfg.Services.Conversion, cfg.Timeouts.Conversion)
	if err != nil {
		log.Error("conversion client failed", "error", err)
		os.Exit(1)
	}
	defer conversion.Close()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.Timeouts.Read,
		IdleTimeout:  cfg.Timeouts.Idle,
		WriteTimeout: cfg.Timeouts.Write,

		Handler: routers.NewAppRouter(&routers.Handlers{
			Currency:   handlers.NewCurrencyHandler(currency, &cfg, log),
			Conversion: handlers.NewConversionHandler(conversion, &cfg, log),
		}, &cfg, log),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Shutdown)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
}
