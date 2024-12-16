package defs

import (
	"context"
	"database/sql"
)

type Repository interface {
	Connect(Config) error
	Close()
	DB() *sql.DB
	SelectContext(context.Context, any, string, ...any) error
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Config interface {
	Validate() error
	DNS() string
	GetDBPath() string
	GetInMemory() bool
}
