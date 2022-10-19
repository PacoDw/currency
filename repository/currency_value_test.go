package repository_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/PacoDw/currency/repository"
	"github.com/jackc/fake"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestBulkInsertCurrencyValue(t *testing.T) {
	TestEnvDBConnectionVariables(t)

	conn := repository.NewSQLConnection(config)

	assert.Condition(t, func() (success bool) { return assert.NotNil(t, conn) })

	req := []repository.CurrencyValue{
		{
			fake.CurrencyCode(),
			1,
			cast.ToFloat64(fake.Latitute()),
			time.Date(2022, 10, 17, 17, 23, 34, 123, time.UTC),
		},
		{
			fake.CurrencyCode(),
			1,
			cast.ToFloat64(fake.Latitute()),
			time.Date(2022, 10, 16, 17, 23, 34, 123, time.UTC),
		},
		{
			fake.CurrencyCode(),
			2,
			cast.ToFloat64(fake.Latitute()),
			time.Date(2022, 10, 15, 17, 23, 34, 123, time.UTC),
		},
		{
			fake.CurrencyCode(),
			3,
			cast.ToFloat64(fake.Latitute()),
			time.Date(2022, 10, 13, 17, 23, 34, 123, time.UTC),
		},
		{
			fake.CurrencyCode(),
			3,
			cast.ToFloat64(fake.Latitute()),
			time.Date(2022, 10, 10, 17, 23, 34, 123, time.UTC),
		},
	}

	err := conn.CurrencyValue.BulkInsert(req)

	assert.Condition(t, func() (success bool) { return assert.Nil(t, err) })
}

func TestGetFinitAndFend(t *testing.T) {
	TestEnvDBConnectionVariables(t)

	conn := repository.NewSQLConnection(config)

	assert.Condition(t, func() (success bool) { return assert.NotNil(t, conn) })

	conn.CheckConn().CurrencyValue.GetFinitAndFend()
}

func TestListCurrenciesByDateRange(t *testing.T) {
	TestEnvDBConnectionVariables(t)

	conn := repository.NewSQLConnection(config)

	assert.Condition(t, func() (success bool) { return assert.NotNil(t, conn) })

	init := time.Date(2022, 10, 7, 17, 22, 34, 0, time.UTC)
	end := time.Date(2022, 10, 18, 17, 22, 34, 0, time.UTC)

	res, err := conn.CurrencyValue.ListCurrenciesByDateRange("Cambodia",
		&init,
		&end,
	)

	assert.Condition(t, func() (success bool) { return assert.Nil(t, err) })
	assert.Condition(t, func() (success bool) { return assert.NotNil(t, res) })

	b, _ := json.MarshalIndent(res, "", "	")

	log.Printf("b: %s\n", b)
}
