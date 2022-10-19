package providers

import (
	"context"
	"time"

	"github.com/PacoDw/currency/logger"
)

// Currencier represent the interface for Currency Providers.
type Currencier interface {
	// GetLatestExchangeRates should retrieve the latest exhange rates from
	// Currency Provider.
	GetLatestExchangeRates(ctx context.Context) (*Metadata, []byte, error)

	// SetLogger represenst an easy way to set an specific logger.
	SetLogger(l *logger.Logger)

	// GetTimeoutRequest helps to check the current Timeout of the request.
	GetTimeoutRequest() time.Duration
}
