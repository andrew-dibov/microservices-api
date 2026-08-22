package clients

import (
	"context"
	"microservices-api/pkg/api/currency"
)

func (client *CurrencyClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

func (client *CurrencyClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, client.conf.CurrencyService.HealthTimeout)
	defer cancel()

	_, err := client.grpc.Rate(ctx, &currency.RateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	})

	return err
}

/* --- --- --- */
