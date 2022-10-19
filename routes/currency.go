package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/PacoDw/currency/repository"
)

// QueryParameter represents the required parameter that could be use for some
// other endpoints.
type DateTimeQueryParameter string

const (
	// Finit represents the id of a company.
	Finit DateTimeQueryParameter = DateTimeQueryParameter("finit")

	// Fend represents the country iso, e.g.: us, ur, etc.
	Fend DateTimeQueryParameter = DateTimeQueryParameter("fend")
)

// RouteParameter represents the required route parameter for endopoints.
type RouteParameter string

// Currency represens the number of the currency provider. this paramter is required
const Currency RouteParameter = RouteParameter("currency")

// CurrencyRoute represents the main rout to handle request accepting a route parameter called 'currency'
// which is required and query parameters with time type such as: finit and fend these parameter are not required.
func CurrencyRoute(repo *repository.SQLConnection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			curr  = r.Context().Value(Currency).(string)
			finit = r.Context().Value(Finit).(time.Time)
			fend  = r.Context().Value(Fend).(time.Time)
		)

		// get the data from the repository
		data, err := repo.CheckConn().CurrencyValue.ListCurrenciesByDateRange(curr, &finit, &fend)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		blob, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(blob); err != nil {
			w.WriteHeader(http.StatusNotFound)

			return
		}
	}
}
