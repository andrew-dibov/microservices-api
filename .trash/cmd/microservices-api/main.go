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
	c := configs.NewAppConfig()
	l := slog.New(map[bool]slog.Handler{
		true:  slog.NewJSONHandler(os.Stdout, nil),
		false: slog.NewTextHandler(os.Stdout, nil),
	}[c.App.Prod])

	l.Info("app config",
		"port", c.App.Port,
		"prod", c.App.Prod,
		"history", c.Addr.Hist,
		"currency", c.Addr.Curr,
		"conversion", c.Addr.Conv,
	)

	curr, err := clients.NewCurrClient(c.Addr.Curr, c.Tout.Curr)
	if err != nil {
		l.Error("currency client", "error", err)
		os.Exit(1)
	}
	defer curr.Close()

	conv, err := clients.NewConvClient(c.Addr.Conv, c.Tout.Conv)
	if err != nil {
		l.Error("conversion client", "error", err)
		os.Exit(1)
	}
	defer conv.Close()

	p := registries.NewPromRegistry()

	s := &http.Server{
		Addr:         ":" + c.App.Port,
		ReadTimeout:  c.Tout.Read,
		IdleTimeout:  c.Tout.Idle,
		WriteTimeout: c.Tout.Write,

		Handler: routers.NewAppRouter(&routers.Handlers{
			App:  handlers.NewAppHandler(curr, conv, p, l),
			Curr: handlers.NewCurrHandler(curr, &c, l),
			Conv: handlers.NewConvHandler(conv, &c, l),
		}, &c, l),
	}

	go func() {

		if c.App.Cert != "" && c.App.Key != "" {
			if err := s.ListenAndServeTLS(c.App.Cert, c.App.Key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				l.Error("server failed", "error", err)
				os.Exit(1)
			}
		} else {
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				l.Error("server failed", "error", err)
				os.Exit(1)
			}
		}
	}()

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	ctx, can := context.WithTimeout(context.Background(), c.Tout.Shutdown)
	defer can()

	if err := s.Shutdown(ctx); err != nil {
		l.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
}
