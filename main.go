package main

import (
	"os"

	"github.com/PacoDw/currency/providers"
	"github.com/PacoDw/currency/repository"
	"github.com/PacoDw/currency/routes"
	"github.com/PacoDw/currency/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
)

var (
	serverPort string
)

func init() {
	serverPort = os.Getenv("SERVER_PORT")
}

func main() {
	// start connection with postgres
	repo := repository.NewSQLConnection(repository.DefaultPostgresConfig())

	// create server and pass a configuration
	s := server.New(
		// pass repository with postgres connection
		server.Repository(repo),

		// setting the currency provider into the server to make the proper requests
		server.CurrencyProvider(
			providers.NewFreeCurrencyAPI(providers.DefaultFreeCurrencyAPIConfig()),
		),

		// setting the request interval for the currency provider
		server.CurrencyRequestInterval(os.Getenv("REQUEST_INTERVAL")),

		// setting some default middlewares to handle the server
		server.UseMidlewares(
			middleware.StripSlashes, // match paths with a trailing slash, strip it, and continue routing through the mux
			middleware.Recoverer,    // recover from panics without crashing server
		),

		// if serverPort is empty by default it takes the port 9000
		server.ListenOn(serverPort),
	)

	s.Route("/currencies", func(r chi.Router) {
		// setting a midleware for this route to handle the incomming calls checking the query parameters
		r.Use(server.ValidateDateTimeQueryParametersMiddleware(
			[]routes.DateTimeQueryParameter{
				routes.Finit,
				routes.Fend,
			}),
		)

		// creating the sub route passing a middleware to handle the currency route parameter and then
		// the route controller.
		r.
			With(server.ValidateRouteParametersMiddleware([]routes.RouteParameter{routes.Currency})).
			Get("/{currency}", routes.CurrencyRoute(repo))
	})

	// start the server
	s.Start()
}
