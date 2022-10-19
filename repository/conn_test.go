package repository_test

import (
	"os"
	"testing"

	"github.com/PacoDw/currency/repository"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	config *repository.Config
)

func setDBConnectionEnvVariables() {
	config = &repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func TestEnvDBConnectionVariables(t *testing.T) {
	if os.Getenv("DB_HOST") == "" ||
		os.Getenv("DB_PORT") == "" ||
		os.Getenv("DB_USER") == "" ||
		os.Getenv("DB_PASS") == "" ||
		os.Getenv("DB_NAME") == "" {
		if err := godotenv.Load("../.env"); err != nil {
			t.Fatalf("error loading .env file %+v", err)
		}
	}

	setDBConnectionEnvVariables()
}

func TestConnectionString(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "my_user")
	os.Setenv("DB_PASS", "password")
	os.Setenv("DB_NAME", "database_name")

	expected := "postgres://my_user:password@localhost:5432/database_name?sslmode=disable"
	got := repository.DefaultPostgresConfig().ToString()

	assert.EqualValues(t, expected, got)
}

func TestNewPostgresConn(t *testing.T) {
	TestEnvDBConnectionVariables(t)

	conn := repository.NewPostgresConn(config)

	assert.Condition(t, func() (success bool) { return assert.NotNil(t, conn) })
}
