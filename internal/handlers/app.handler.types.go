package handlers

import (
	"log/slog"
	"microservices-api/internal/clients"
	"microservices-api/internal/registries"
)

type AppHandler struct {
	curr *clients.CurrClient
	conv *clients.ConvClient
	log  *slog.Logger
	preg *registries.PromRegistry
}
