// Package repository contains persistence operations for wallets, transfers,
// ledger entries, and idempotency records.
//
// Store only wires the connection pool; the query/locking logic that
// implements the transfer design lives here as it's written.
package repository

import "github.com/jackc/pgx/v5/pgxpool"

// Store wraps the shared database connection pool used by repository methods.
type Store struct {
	Pool *pgxpool.Pool
}

// NewStore constructs a Store around an existing pgx connection pool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}
