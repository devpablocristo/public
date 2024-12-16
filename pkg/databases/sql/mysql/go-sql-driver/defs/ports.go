package defs

import "database/sql"

type Repository interface {
	Close()
	DB() *sql.DB
}
