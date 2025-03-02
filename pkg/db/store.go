package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *SqlStore

type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) *SqlStore {
	return &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
