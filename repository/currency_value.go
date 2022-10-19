package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// CurrencyValueRepository defines the interface that device must satisfy.
type CurrencyValueRepository interface {
	BulkInsert(cvs []CurrencyValue) error
	ListCurrenciesByDateRange(currency string, finit, fend *time.Time) ([]CurrencyValue, error)
	GetFinitAndFend() (finit, fend time.Time, err error)
}

// CurrencyValueSQLService represents a sqlService type.
type CurrencyValueSQLService sqlService

// CurrencyValueSQLService validate if it satisfy the own interface, that means
// that all sqlService can be implement its own interface but it must be
// a sqlService type.
var _ CurrencyValueRepository = &CurrencyValueSQLService{}

type CurrencyValue struct {
	Name         string    `json:"name,omitempty"`
	RequestID    int64     `json:"request_id,omitempty"`
	Value        float64   `json:"value,omitempty"`
	LastUdatedAt time.Time `json:"last_updated_at,omitempty"`
}

// BulkInsert inserts all currencies values from the Currency provider into the database.
func (service *CurrencyValueSQLService) BulkInsert(cvs []CurrencyValue) error {
	var (
		placeholders = []string{}
		vals         = []interface{}{}
	)

	for i := range cvs {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)",
			i*4+1,
			i*4+2,
			i*4+3,
			i*4+4,
		))

		vals = append(vals, cvs[i].Name, cvs[i].RequestID, cvs[i].Value, cvs[i].LastUdatedAt)
	}

	tx, err := service.db.Begin()
	if err != nil {
		return errors.Wrap(err, "could not start a new transaction")
	}

	insertStatement := fmt.Sprintf(`
		INSERT INTO currencies_values
			(
				name,
				request_id,
				value,
				last_updated_at
			)
		VALUES 
			%s
	;`, strings.Join(placeholders, ","))

	_, err = tx.Exec(insertStatement, vals...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to make a rollback")
		}

		return errors.Wrap(err, "failed to insert multiple records at once")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// ListCurrenciesByDateRange represents a function to retrieve data with the next parameters:
// currency => it must be 3 letters and it could be 'all' as a value, it not accepts numbers, it is required
// finit    => is a start date is optional
// fend     => is an end date is optional
// Note: these paramters are filters.
func (service *CurrencyValueSQLService) ListCurrenciesByDateRange(currency string, finit, fend *time.Time) ([]CurrencyValue, error) {
	if finit == nil || fend == nil || finit.IsZero() || fend.IsZero() {
		init, end, err := service.GetFinitAndFend()
		if err != nil {
			return nil, err
		}

		if finit == nil || finit.IsZero() {
			finit = &init
		}

		if fend == nil || fend.IsZero() {
			fend = &end
		}
	}

	var stmCond string
	if !strings.EqualFold(currency, "all") {
		stmCond = fmt.Sprintf("AND name = '%s'", currency)
	}

	rows, err := service.db.Query(fmt.Sprintf(`
		SELECT 
			name,
			request_id,
			value,
			last_updated_at::TIMESTAMP
		FROM
			currencies_values
		WHERE 
			last_updated_at::TIMESTAMP >= '%s'::TIMESTAMP 
		AND 
			last_updated_at::TIMESTAMP <='%s'::TIMESTAMP
		%s;
`, finit.Format("2006-01-02 15:04:05"), fend.Format("2006-01-02 15:04:05"), stmCond))
	if err != nil {
		return nil, errors.Wrap(err, "failed to range between dates")
	}

	vals := make([]CurrencyValue, 0)

	for rows.Next() {
		var cv CurrencyValue

		if err := rows.Scan(
			&cv.Name,
			&cv.RequestID,
			&cv.Value,
			&cv.LastUdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scanning multiple records")
		}

		vals = append(vals, cv)
	}

	return vals, nil
}

// GetFinitAndFend gets the first date and the last date inserted in the database.
func (service *CurrencyValueSQLService) GetFinitAndFend() (finit, fend time.Time, err error) {
	finit = time.Time{}
	fend = time.Time{}

	rows, err := service.db.Query(`
		SELECT last_updated_at::TIMESTAMP
		FROM (
			SELECT last_updated_at::TIMESTAMP,
				ROW_NUMBER() OVER (ORDER BY last_updated_at::TIMESTAMP DESC) AS rn,
				COUNT(*) OVER () AS total_count
			FROM currencies_values
		) t
		WHERE rn = 1
		OR rn = total_count
		ORDER BY last_updated_at::TIMESTAMP DESC;
	`)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "failed to get the first and the last date of the table")
	}

	rows.Next()

	if err := rows.Scan(&fend); err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "failed to scan attributes")
	}

	rows.Next()

	if err := rows.Scan(&finit); err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "failed to scan attributes")
	}

	if err := rows.Err(); err != nil {
		return time.Time{}, time.Time{}, err
	}

	return
}
