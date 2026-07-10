package main

import (
	"microservices-api/internal/clients"
	"microservices-api/internal/configs"
	"microservices-api/internal/handlers"
	"microservices-api/internal/registries"
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

	log := slog.New(map[bool]slog.Handler{
		true:  slog.NewJSONHandler(os.Stdout, nil),
		false: slog.NewTextHandler(os.Stdout, nil),
	}[cfg.Prod])

	log.Info("app config",
		"port", cfg.Port,
		"prod", cfg.Prod,
		"history", cfg.Services.Hist,
		"currency", cfg.Services.Curr,
		"conversion", cfg.Services.Conv,
	)

	curr, err := clients.NewCurrClient(cfg.Services.Curr, cfg.Timeouts.Curr)
	if err != nil {
		log.Error("currency client",
			"error", err,
		)
		os.Exit(1)
	}
	defer curr.Close()

	conv, err := clients.NewConvClient(cfg.Services.Conv, cfg.Timeouts.Conv)
	if err != nil {
		log.Error("conversion client",
			"error", err,
		)
		os.Exit(1)
	}
	defer conv.Close()

	preg := registries.NewPromRegistry()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.Timeouts.Read,
		IdleTimeout:  cfg.Timeouts.Idle,
		WriteTimeout: cfg.Timeouts.Write,

		Handler: routers.NewAppRouter(&routers.Handlers{
			App:  handlers.NewAppHandler(curr, conv, preg, log),
			Curr: handlers.NewCurrHandler(curr, &cfg, log),
			Conv: handlers.NewConvHandler(conv, &cfg, log),
		}, &cfg, log),
	}

	go func() {

		if cfg.Cert != "" && cfg.Key != "" {
			if err := srv.ListenAndServeTLS(cfg.Cert, cfg.Key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("server failed",
					"error", err,
				)
				os.Exit(1)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("server failed",
					"error", err,
				)
				os.Exit(1)
			}
		}
	}()

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	ctx, can := context.WithTimeout(context.Background(), cfg.Timeouts.Shutdown)
	defer can()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown failed",
			"error", err,
		)
		os.Exit(1)
	}
}
