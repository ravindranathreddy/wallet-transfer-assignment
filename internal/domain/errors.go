// Package domain holds the wallet transfer entities, state machine, and
// validation rules. Intentionally left for the core design/implementation
// (entities, transfer state transitions, invariants) rather than scaffolding.
package domain

import "errors"

var (
	ErrSameWallet = errors.New("cannot transfer to the same wallet")

	ErrInvalidAmount = errors.New("transfer amount must be greater than zero")

	ErrInvalidTransition = errors.New("invalid state transition")

	ErrIdempotencyKeyRequired = errors.New("idempotencyKey is required")
	ErrWalletIdRequired       = errors.New("WalletId cannot be empty")
)
