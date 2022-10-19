package repository_test

import (
	"testing"
	"time"

	"github.com/PacoDw/currency/repository"
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
)

func TestInsertRequestStatus(t *testing.T) {
	TestEnvDBConnectionVariables(t)

	conn := repository.NewSQLConnection(config)

	assert.Condition(t, func() (success bool) { return assert.NotNil(t, conn) })

	reqs := []repository.RequestStatus{
		{
			time.Duration(1 * time.Second).String(),
			fake.DomainName(),
			fake.Gender(),
			time.Now(),
		},
		{
			time.Duration(1 * time.Second).String(),
			fake.DomainName(),
			fake.Gender(),
			time.Now(),
		},
		{
			time.Duration(1 * time.Second).String(),
			fake.DomainName(),
			fake.Gender(),
			time.Now(),
		},
	}

	for i := range reqs {
		reqID, err := conn.RequestStatus.Insert(reqs[i])

		assert.Condition(t, func() (success bool) { return assert.Nil(t, err) })
		assert.EqualValues(t, i+1, reqID)
	}
}
