package mysql

import (
	"context"
	"database/sql"

	"github.com/ytake/protoactor-go-cqrs-example/internal/config"
)

// NewConn makes a new connection to the database
func NewConn(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MySQLMaxOpenConns)
	db.SetMaxIdleConns(config.MySQLMaxIdleConns)
	db.SetConnMaxLifetime(config.MySQLConnMaxLifetime)
	return db, nil
}

// RegistrationUserExecutor is Query Executor
type RegistrationUserExecutor interface {
	AddUserIfNotExists(ctx context.Context, param AddUserParams) error
}

type RegistrationUserQueryExecutor interface {
	GetRegistrationUsers(ctx context.Context) ([]User, error)
}
