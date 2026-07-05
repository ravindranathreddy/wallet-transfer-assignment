//go:build integration

package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ravindranathreddy/wallet-transfer-assignment/internal/domain"
	"github.com/ravindranathreddy/wallet-transfer-assignment/internal/repository"
	"github.com/ravindranathreddy/wallet-transfer-assignment/internal/testutil"
)

func TestTransferRepository_InsertAndGet(t *testing.T) {
	pool := testutil.NewPool(t)
	testutil.TruncateAll(t, pool)
	testutil.SeedWallet(t, pool, "wallet_1", 100)
	testutil.SeedWallet(t, pool, "wallet_2", 0)

	repo := repository.NewTransferRepository()
	transfer, err := domain.NewPendingTransfer(uuid.New(), "wallet_1", "wallet_2", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.InsertTransfer(context.Background(), pool, transfer); err != nil {
		t.Fatalf("insert transfer: %v", err)
	}

	got, err := repo.GetTransferByID(context.Background(), pool, transfer.ID)
	if err != nil {
		t.Fatalf("get transfer: %v", err)
	}
	if got.Status != domain.StatusPending {
		t.Errorf("expected status PENDING, got %s", got.Status)
	}
	if got.Amount != 50 {
		t.Errorf("expected amount 50, got %d", got.Amount)
	}
	if got.FromWalletID != "wallet_1" || got.ToWalletID != "wallet_2" {
		t.Errorf("expected wallet_1->wallet_2, got %s->%s", got.FromWalletID, got.ToWalletID)
	}
}

func TestTransferRepository_GetTransferByID_NotFound(t *testing.T) {
	pool := testutil.NewPool(t)
	testutil.TruncateAll(t, pool)

	repo := repository.NewTransferRepository()
	_, err := repo.GetTransferByID(context.Background(), pool, uuid.New())
	if !errors.Is(err, domain.ErrTransferNotFound) {
		t.Errorf("expected ErrTransferNotFound, got %v", err)
	}
}

func TestTransferRepository_UpdateTransferStatus(t *testing.T) {
	pool := testutil.NewPool(t)
	testutil.TruncateAll(t, pool)
	testutil.SeedWallet(t, pool, "wallet_1", 100)
	testutil.SeedWallet(t, pool, "wallet_2", 0)

	repo := repository.NewTransferRepository()
	transfer, err := domain.NewPendingTransfer(uuid.New(), "wallet_1", "wallet_2", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.InsertTransfer(context.Background(), pool, transfer); err != nil {
		t.Fatalf("insert transfer: %v", err)
	}

	reason := "insufficient balance"
	if err := repo.UpdateTransferStatus(context.Background(), pool, transfer.ID, domain.StatusFailed, &reason); err != nil {
		t.Fatalf("update status: %v", err)
	}

	got, err := repo.GetTransferByID(context.Background(), pool, transfer.ID)
	if err != nil {
		t.Fatalf("get transfer: %v", err)
	}
	if got.Status != domain.StatusFailed {
		t.Errorf("expected status FAILED, got %s", got.Status)
	}
	if got.FailureReason == nil || *got.FailureReason != reason {
		t.Errorf("expected failure reason %q, got %v", reason, got.FailureReason)
	}
}
