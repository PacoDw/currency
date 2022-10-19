package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PacoDw/currency/providers"
	"github.com/PacoDw/currency/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// FuncOptionType represents the option type in number which the order of the
// options tell us which needs to be applied first.
type FuncOptionType int

const (
	CURRENCYPROVIDER FuncOptionType = iota
	CURRENCYREQUESTINTERVAL
	LISTENON
	LOGGER
	MIDLEWARES
	ROUTES
	HANDLER
)

// Option represents an Option interface that can be set in the server constructor.
type Option interface {
	Apply(*Server)
	Key() FuncOptionType
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc struct {
	key      FuncOptionType
	callback func(*Server)
}

// Apply applies the current option in the server.
func (f optionFunc) Apply(s *Server) {
	f.callback(s)
}

// Key represents the number of the options this helps to check what options
// needs to be applied first.
func (f optionFunc) Key() FuncOptionType {
	return f.key
}

// Repository allows to set a provider.
func Repository(repo *repository.SQLConnection) Option {
	if repo == nil {
		panic("the currency option must not be nil")
	}

	return optionFunc{
		key: CURRENCYPROVIDER,
		callback: func(s *Server) {
			s.repo = repo
		},
	}
}

// CurrencyProvider allows to set a provider.
func CurrencyProvider(currency providers.Currencier) Option {
	if currency == nil {
		panic("the currency option must not be nil")
	}

	return optionFunc{
		key: CURRENCYPROVIDER,
		callback: func(s *Server) {
			s.currencyClient = currency

			s.currencyClient.SetLogger(s.logger)
		},
	}
}

// CurrencyRequestInterval allows to set an interval of time to call the provider.
func CurrencyRequestInterval(requestInterval string) Option {
	return optionFunc{
		key: CURRENCYPROVIDER,
		callback: func(s *Server) {
			d, err := time.ParseDuration(requestInterval)
			if err != nil {
				s.logger.Warn("the requestInterval is not correct using default value (10s)",
					zap.String("requestInterval", requestInterval),
					zap.String("err", err.Error()),
				)

				return
			}

			s.currencyRequestInterval = d
		},
	}
}

// ListenOn optionally specifies the TCP address for the server to listen on,
// in the form "host:port". If empty, ":http" (port 9000) is used.
// The service names are defined in RFC 6335 and assigned by IANA.
// See net.Dial for details of the address format.
func ListenOn(addr string) Option {
	return optionFunc{
		key: LISTENON,
		callback: func(s *Server) {
			if addr == "" {
				addr = "9000"
			}

			s.Addr = fmt.Sprintf(":%s", addr)
		},
	}
}

// UseMidlewares allows to set different middlewares.
func UseMidlewares(middlewares ...func(http.Handler) http.Handler) Option {
	return optionFunc{
		key: MIDLEWARES,
		callback: func(s *Server) {
			s.Use(middlewares...)
		},
	}
}

// Routes allows to set different routes in one pattern name.
func Routes(pattern string, callback func(r chi.Router)) Option {
	return optionFunc{
		key: ROUTES,
		callback: func(s *Server) {
			s.Route(pattern, callback)
		},
	}
}

// Handler allows to set a one route and a handler.
func Handler(path string, handler http.Handler) Option {
	return optionFunc{
		key: HANDLER,
		callback: func(s *Server) {
			s.Handle(path, handler)
		},
	}
}
