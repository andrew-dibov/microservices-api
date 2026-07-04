package routers

import "microservices-api/internal/handlers"

type Handlers struct {
	Currency   *handlers.CurrencyHandler
	Conversion *handlers.ConversionHandler
}
