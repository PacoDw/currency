package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbConn *sql.DB
)

// Config is used to set the database configuration to connect to.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// ToString parses the config to a readable connection string.
func (c *Config) ToString() string {
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Path:     c.DBName,
		RawQuery: "sslmode=disable",
	}

	return u.String()
}

// DefaultPostgresConfig helps to set the env variables into a the config struct.
func DefaultPostgresConfig() *Config {
	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

// NewPostgresConn accepts a database configuration to create a new mssql connection.
func NewPostgresConn(c *Config) *sql.DB {
	if dbConn != nil {
		if err := dbConn.Ping(); err != nil {
			dbConn = nil
		}
	}

	if dbConn == nil {
		var err error

		dbConn, err = sql.Open("postgres", c.ToString())
		if err != nil {
			log.Panicf("can't reach database: %s", err)
		}
	}

	return dbConn
}

// NewSQLConnection creates a new SQLConnection with all repositories inside,
// internally it connects with the database.
func NewSQLConnection(c *Config) *SQLConnection {
	return New(NewPostgresConn(c), c)
}
