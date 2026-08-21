package clients

import (
	"context"
	"microservices-api/pkg/api/conversion"
)

func (cl *ConversionClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}
	return nil
}

func (cl *ConversionClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, cl.conf.ConversionService.HealthTimeout)
	defer cancel()

	_, err := cl.grpc.Convert(ctx, &conversion.ConvertRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       1,
	})

	return err
}

/* --- --- --- */

func (cl *ConversionClient) Convert(ctx context.Context, fromCurrency string, toCurrency string, amount float64) (*conversion.ConvertResponse, error) {
	// REQ ID MIDDLEWARE

	ctx, cancel := context.WithTimeout(ctx, cl.conf.ConversionService.Timeout)
	defer cancel()

	return cl.grpc.Convert(ctx, &conversion.ConvertRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
	})
}
