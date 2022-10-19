package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/PacoDw/currency/providers"
	"github.com/PacoDw/currency/repository"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
}

func TestProviderJob(t *testing.T) {
	repo := repository.NewSQLConnection(repository.DefaultPostgresConfig())
	currencyProvider := providers.NewFreeCurrencyAPI(providers.DefaultFreeCurrencyAPIConfig())

	s := New(
		Repository(repo),
		CurrencyProvider(currencyProvider),
		CurrencyRequestInterval(os.Getenv("REQUEST_INTERVAL")),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go s.RunProviderJob(ctx)

	<-ctx.Done()
}
