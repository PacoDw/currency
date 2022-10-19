package providers_test

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/PacoDw/currency/providers"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("FREE_CURRENCY_URL", "https://api.currencyapi.com/v3/latest")
	os.Setenv("FREE_CURRENCY_API_KEY", "API_KEY")
	os.Setenv("REQUEST_TIMEOUT", "3s")
}

func TestCurrencyConfig(t *testing.T) {
	u, err := url.Parse("https://api.currencyapi.com/v3/latest")
	assert.EqualValues(t, err, nil)

	t.Run("empty APIKey", func(t *testing.T) {
		cc := &providers.CurrencyConfig{
			URL:     u,
			APIKey:  "",
			Timeout: time.Duration(3 * time.Second),
		}

		err = cc.Valid()

		assert.EqualError(t, err, "the APIKey attribute must not be empty")
	})

	t.Run("nil url", func(t *testing.T) {
		cc := &providers.CurrencyConfig{
			URL:     nil,
			APIKey:  "API_KEY",
			Timeout: time.Duration(3 * time.Second),
		}

		err := cc.Valid()

		assert.EqualError(t, err, "the URL attribute must not be empty")
	})

	t.Run("empty timeout", func(t *testing.T) {
		cc := &providers.CurrencyConfig{
			URL:     u,
			APIKey:  "API_KEY",
			Timeout: 0,
		}

		err := cc.Valid()

		assert.EqualError(t, err, "the Timeout attribute must not be 0")
	})

	t.Run("Success", func(t *testing.T) {
		expected := &providers.CurrencyConfig{
			URL:     u,
			APIKey:  "API_KEY",
			Timeout: time.Duration(3 * time.Second),
		}

		got := providers.DefaultFreeCurrencyAPIConfig()

		assert.EqualValues(t, expected, got)
	})
}
