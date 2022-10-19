package providers

import (
	"errors"
	"net/url"
	"time"

	"github.com/PacoDw/currency/logger"
)

// CurrencyConfig represents the needed configuration to create a Currency Provider.
type CurrencyConfig struct {
	// URL represents the main url of the API call, this is a required attribute.
	URL *url.URL

	// APIKey represents the API KEY of the Currency Provider, this is a required
	// attribute.
	APIKey string

	// Timeout represents the timeout of each request made by the Currency Provider
	// this is a required attribute.
	Timeout time.Duration

	// Logger is an optional attribute, plz use checkLogger to set a default logger
	Logger *logger.Logger
}

// checkLogger checks if there is an logger set, if not it will set a logger with
// a default configuration.
func (cfg *CurrencyConfig) checkLogger() {
	if cfg.Logger != nil {
		return
	}

	cfg.Logger = logger.NewLogger(logger.DefaultEnvLoggerConfig())
}

// SetLogger internally checks as first step if the l parameter is nil with
// checkLogger function if so then it will set a default logger with a default
// cnfiguration. Otherwise, it will set the logger using l parameter.
func (cfg *CurrencyConfig) SetLogger(l *logger.Logger) {
	if l == nil {
		cfg.checkLogger()
	}

	cfg.Logger = l
}

// Valid validates if the required attributes has been set, if not it will return
// and error specifying what parameter is empty o nil.
func (cfg *CurrencyConfig) Valid() error {
	if cfg.APIKey == "" {
		return errors.New("the APIKey attribute must not be empty")
	}

	if cfg.URL == nil {
		return errors.New("the URL attribute must not be empty")
	}

	if cfg.Timeout == 0 {
		return errors.New("the Timeout attribute must not be 0")
	}

	return nil
}
