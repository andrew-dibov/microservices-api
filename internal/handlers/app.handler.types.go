package handlers

import (
	"log/slog"
	"microservices-api/internal/clients"
)

type AppHandler struct {
	curr *clients.CurrClient
	conv *clients.ConvClient
	log  *slog.Logger
}
