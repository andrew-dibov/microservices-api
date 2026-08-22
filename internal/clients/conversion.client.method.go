package clients

import (
	"context"
	"microservices-api/pkg/api/conversion"
)

func (client *ConversionClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

func (client *ConversionClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, client.conf.ConversionService.HealthTimeout)
	defer cancel()

	_, err := client.grpc.Convert(ctx, &conversion.ConvertRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       1,
	})

	return err
}

/* --- --- --- */

func (client *ConversionClient) Convert(ctx context.Context, fromCurrency string, toCurrency string, amount float64) (*conversion.ConvertResponse, error) {
	// REQ ID MIDDLEWARE

	ctx, cancel := context.WithTimeout(ctx, client.conf.ConversionService.Timeout)
	defer cancel()

	return client.grpc.Convert(ctx, &conversion.ConvertRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
	})
}
