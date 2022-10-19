package repository

import "database/sql"

// nolint // sqlService represents a type for each service created, so all the servies like
// device in device.go must implement this type.
type sqlService struct {
	db     *sql.DB
	config *Config
}

// SQLConnection represenst the SQLConnection that contains all the services created.
// Note: If you has been created a new service it must be listed in this struct.
type SQLConnection struct {
	sqlService    *sqlService
	RequestStatus RequestStatusRepository
	CurrencyValue CurrencyValueRepository
}

// CheckConn check if the connection is alive if not then it will connect again.
func (conn *SQLConnection) CheckConn() *SQLConnection {
	conn.sqlService.db = NewPostgresConn(conn.sqlService.config)

	return conn
}

// New creates a new SQLConnection with all the services in it
// Note: If you has been created a new service it must be listed in this struct.
func New(db *sql.DB, cfg *Config) *SQLConnection {
	sqls := &sqlService{
		db,
		cfg,
	}

	return &SQLConnection{
		sqlService:    sqls,
		RequestStatus: (*RequestStatusSQLService)(sqls),
		CurrencyValue: (*CurrencyValueSQLService)(sqls),
	}
}
