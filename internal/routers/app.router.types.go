package routers

import "microservices-api/internal/handlers"

type Handlers struct {
	Curr *handlers.CurrHandler
	Conv *handlers.ConvHandler
}
