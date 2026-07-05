package routers

import "microservices-api/internal/handlers"

type Handlers struct {
	App  *handlers.AppHandler
	Curr *handlers.CurrHandler
	Conv *handlers.ConvHandler
}
