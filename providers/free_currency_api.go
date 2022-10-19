package providers

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

// freeCurrencyApi is a Currency Provider.
type freeCurrencyApi struct {
	*CurrencyConfig

	authHeader *http.Header
	client     *http.Client
}

// Metadata helps to register all stats about the request.
type Metadata struct {
	Elapsed     string
	URL         string
	Status      string
	RequestedAt time.Time
	Error       error
}

// GetLatestExchangeRates returns the stats as metadata, the bytes represents the result of the request, and
// finally an error if the case.
func (fc *freeCurrencyApi) GetLatestExchangeRates(ctx context.Context) (*Metadata, []byte, error) {
	var (
		cancel context.CancelFunc
		resCh  = make(chan []byte)
		errCh  = make(chan error)
	)

	// creates the metadata struct
	meta := &Metadata{
		URL:         fc.URL.String(),
		RequestedAt: time.Now(),
	}

	// when the process finishs we will register the elapsed time
	defer func() { meta.Elapsed = time.Since(meta.RequestedAt).String() }()

	// set the timeout of each request
	ctx, cancel = context.WithTimeout(ctx, fc.Timeout)
	defer cancel()

	go func() {
		// preparing the request
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fc.URL.String(), http.NoBody)
		if err != nil {
			errCh <- err

			return
		}

		// getting the apikey
		req.Header = *fc.authHeader

		// making the request
		res, err := fc.client.Do(req)
		if err != nil {
			errCh <- err

			return
		}

		blob, err := io.ReadAll(res.Body)
		if err != nil {
			errCh <- err

			return
		}
		defer res.Body.Close()

		resCh <- blob
	}()

	select {
	case <-ctx.Done():
		meta.Error = errors.Wrap(ctx.Err(), "the proccess ends before finish by timeout")
		meta.Status = "failure"

		return meta, nil, meta.Error
	case res := <-resCh:
		meta.Status = "success"

		return meta, res, nil
	case err := <-errCh:
		meta.Error = err
		meta.Status = "failure"

		return meta, nil, err
	}
}

// GetTimeoutRequest returns the timeout set for requests.
func (fc *freeCurrencyApi) GetTimeoutRequest() time.Duration {
	return fc.Timeout
}

// NewFreeCurrencyAPI creates a new Free Currency Api Provider.
func NewFreeCurrencyAPI(cfg *CurrencyConfig) Currencier {
	if cfg == nil {
		panic("the CurrencyConfig must not be nil")
	}

	if err := cfg.Valid(); err != nil {
		panic(err)
	}

	return &freeCurrencyApi{
		CurrencyConfig: cfg,
		authHeader:     &http.Header{"apikey": []string{cfg.APIKey}},
		client:         &http.Client{},
	}
}

// DefaultFreeCurrencyAPIConfig sets the default configuration taking the proper env variables.
func DefaultFreeCurrencyAPIConfig() *CurrencyConfig {
	u, err := url.Parse(os.Getenv("FREE_CURRENCY_URL"))
	if err != nil {
		panic(err)
	}

	reqTimeout, err := time.ParseDuration(os.Getenv("REQUEST_TIMEOUT"))
	if err != nil {
		panic(err)
	}

	return &CurrencyConfig{
		URL:     u,
		APIKey:  os.Getenv("FREE_CURRENCY_API_KEY"),
		Timeout: reqTimeout,
	}
}
