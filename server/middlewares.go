package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PacoDw/currency/routes"
	"github.com/go-chi/chi/v5"
)

// ValidateRouteParametersMiddleware validates that the incoming request has the proper route parameters
// if not it is descarted.
func ValidateRouteParametersMiddleware(rps []routes.RouteParameter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Validate the routes parameters are passing correctly
			for i := range rps {
				rp := chi.URLParam(r, string(rps[i]))

				// check if the route parameter is empty
				if rp == "" {
					w.WriteHeader(http.StatusBadRequest)

					w.Write([]byte(fmt.Sprintf(`{"error":"the route parameter is empty (%s)"}`, rps[i])))

					return
				}

				// check if the route parameter not contains 3 letters
				if len(rp) != 3 {
					w.WriteHeader(http.StatusBadRequest)

					w.Write([]byte(fmt.Sprintf(`{"error":"bad route parameter (%s) with value (%s). it must contain only 3 letters"}`, rps[i], rp)))

					return
				}

				// check if the route parameter contains any number
				containsNumber, err := regexp.MatchString("[0-9]+", string(rp))
				if containsNumber || err != nil {
					w.WriteHeader(http.StatusBadRequest)

					w.Write([]byte(fmt.Sprintf(`{"error":"bad route parameter (%s) with value (%s). it must contains a number or is invalid"}`, rps[i], rp)))

					return
				}

				// save the current route parameter in upper case
				ctx = context.WithValue(ctx, rps[i], strings.ToUpper(rp))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// ValidateDateTimeQueryParametersMiddleware validates that the incoming request has the proper query parameters
// if not it is descarted.
func ValidateDateTimeQueryParametersMiddleware(qrps []routes.DateTimeQueryParameter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Validate that all query parameters that are required are passed correctly
			for i := range qrps {
				// getting the current query parameter
				v := r.URL.Query().Get(string(qrps[i]))

				// converting the current query parameter in time with the specific format
				t, err := time.Parse("2006-01-02T15:04:05", v)
				if err != nil && v != "" {
					w.WriteHeader(http.StatusBadRequest)

					w.Write([]byte(fmt.Sprintf(`{"error":"bad query parameter %s with value %s"}`, qrps[i], v)))

					return
				}

				// save the current query parameter
				ctx = context.WithValue(ctx, qrps[i], t)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
