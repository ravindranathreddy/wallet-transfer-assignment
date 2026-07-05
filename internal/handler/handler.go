// Package handler maps HTTP requests to service calls. Handlers stay thin:
// request validation and transport mapping only, no business logic.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ravindranathreddy/wallet-transfer-assignment/internal/domain"
	"github.com/ravindranathreddy/wallet-transfer-assignment/internal/service"
)

// TransferHandler serves the /transfers endpoints.
type TransferHandler struct {
	service *service.TransferService
}

// NewTransferHandler constructs a TransferHandler around a TransferService.
func NewTransferHandler(svc *service.TransferService) *TransferHandler {
	return &TransferHandler{service: svc}
}

type createTransferRequestDTO struct {
	IdempotencyKey string `json:"idempotencyKey"`
	FromWalletID   string `json:"fromWalletId"`
	ToWalletID     string `json:"toWalletId"`
	Amount         int64  `json:"amount"`
}

type transferResponseDTO struct {
	TransferID    string  `json:"transferId"`
	Status        string  `json:"status"`
	FromWalletID  string  `json:"fromWalletId"`
	ToWalletID    string  `json:"toWalletId"`
	Amount        int64   `json:"amount"`
	FailureReason *string `json:"failureReason"`
	CreatedAt     string  `json:"createdAt"`
}

type errorResponseDTO struct {
	Error string `json:"error"`
}

// CreateTransfer handles POST /transfers.
// TODO: request validation, invoking TransferService, response/error mapping.
func (h *TransferHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var req createTransferRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	if req.IdempotencyKey == "" {
		writeError(w, http.StatusBadRequest, domain.ErrIdempotencyKeyRequired.Error())
		return
	}
	if req.FromWalletID == "" || req.ToWalletID == "" {
		writeError(w, http.StatusBadRequest, domain.ErrWalletIdRequired.Error())
		return
	}
	if req.FromWalletID == req.ToWalletID {
		writeError(w, http.StatusBadRequest, domain.ErrSameWallet.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, domain.ErrInvalidAmount.Error())
		return
	}
	writeError(w, http.StatusNotImplemented, "not implemented")
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponseDTO{Error: message})
}
