package repository

import (
	"time"

	"github.com/pkg/errors"
)

// RequestStatusRepository defines the interface that device must satisfy.
type RequestStatusRepository interface {
	Insert(rs RequestStatus) (int64, error)
}

// RequestStatusSQLService represents a sqlService type.
type RequestStatusSQLService sqlService

// RequestStatusSQLService validate if it satisfy the own interface, that means
// that all sqlService can be implement its own interface but it must be
// a sqlService type.
var _ RequestStatusRepository = &RequestStatusSQLService{}

type RequestStatus struct {
	TimeElapsed string
	URL         string
	Status      string
	RequestedAt time.Time
}

// Insert creates registers into the database about all request made.
func (service *RequestStatusSQLService) Insert(rs RequestStatus) (int64, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "could not start a new transaction")
	}

	res := tx.QueryRow(`
		INSERT INTO requests_status
			(
				time_elapsed,
				url,
				status,
				requested_at
			)
		VALUES 
			($1,$2,$3,$4)
		RETURNING id;`, rs.TimeElapsed, rs.URL, rs.Status, rs.RequestedAt)
	if res.Err() != nil {
		if err := tx.Rollback(); err != nil {
			return 0, errors.Wrap(err, "failed to make a rollback")
		}

		return 0, errors.Wrap(err, "failed to insert a record")
	}

	var id int64
	if err := res.Scan(&id); err != nil {
		return 0, errors.Wrap(err, "failed to scan attributes")
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "failed to commit transaction")
	}

	return id, nil
}
