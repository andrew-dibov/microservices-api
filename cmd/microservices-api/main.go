package main

import (
	"microservices-api/internal/configs"
	"microservices-api/internal/routers"

	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := configs.NewAppConfig()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info("app config loaded",
		"port", cfg.Port, "history", cfg.Services.History, "currency", cfg.Services.Currency, "conversion", cfg.Services.Conversion)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.Timeouts.Read,
		IdleTimeout:  cfg.Timeouts.Idle,
		WriteTimeout: cfg.Timeouts.Write,

		Handler: routers.NewAppRouter(&cfg, log),
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
		os.Exit(1)
	}
}
