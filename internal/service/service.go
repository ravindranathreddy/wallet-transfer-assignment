// Package service contains transfer business logic: idempotency handling,
// transfer orchestration, and ledger/balance invariants.
//
// TransferService only wires its repository dependency; the orchestration
// logic that implements the transfer design lives here as it's written.
package service

import "github.com/ravindranathreddy/wallet-transfer-assignment/internal/repository"

// TransferService orchestrates wallet-to-wallet transfers.
type TransferService struct {
	store *repository.Store
}

// NewTransferService constructs a TransferService around a repository Store.
func NewTransferService(store *repository.Store) *TransferService {
	return &TransferService{store: store}
}
