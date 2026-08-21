package clients

import (
	"context"
	"microservices-api/pkg/api/currency"
)

func (cl *CurrencyClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}
	return nil
}

func (cl *CurrencyClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, cl.conf.CurrencyService.HealthTimeout)
	defer cancel()

	_, err := cl.grpc.Rate(ctx, &currency.RateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	})

	return err
}

/* --- --- --- */
