package main

import (
	"context"
	"errors"
	"log/slog"
	"microservices-api/internal/configs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := configs.NewAppConfig()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info("config loaded",
		"port", cfg.Port,
		"history", cfg.Services.History,
		"currency", cfg.Services.Currency,
		"conversion", cfg.Services.Conversion,
	)

	/* --- */

	srv := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  cfg.Timeouts.Read,
		IdleTimeout:  cfg.Timeouts.Idle,
		WriteTimeout: cfg.Timeouts.Write,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "error", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("server stopping")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Shutdown)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server killed", "error", err.Error())
		os.Exit(1)
	}

	log.Info("server stopped")
}
