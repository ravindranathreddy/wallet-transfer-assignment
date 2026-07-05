// Package handler maps HTTP requests to service calls. Handlers stay thin:
// request validation and transport mapping only, no business logic.
package handler

import (
	"net/http"

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

// CreateTransfer handles POST /transfers.
// TODO: request validation, invoking TransferService, response/error mapping.
func (h *TransferHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
