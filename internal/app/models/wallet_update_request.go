package models

import (
	"errors"

	"github.com/google/uuid"
)

const (
	OperationDeposit  = "DEPOSIT"
	OperationWithdraw = "WITHDRAW"
	maxAmount         = 1000000
)

var (
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrInvalidAmount        = errors.New("amount must be greater than 0")
	ErrExceedMaxAmount      = errors.New("amount exceed maximum")
	ErrInvalidWalletID      = errors.New("invalid wallet ID")
)

// WalletUpdateRequest represents request to update wallet
type WalletUpdateRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

// Validate validates wallet update request
func (r *WalletUpdateRequest) Validate() error {
	if r.WalletID == uuid.Nil {
		return ErrInvalidWalletID
	}
	if r.Amount <= 0 {
		return ErrInvalidAmount
	}
	if r.Amount > maxAmount {
		return ErrExceedMaxAmount
	}
	if r.OperationType != OperationDeposit && r.OperationType != OperationWithdraw {
		return ErrInvalidOperationType
	}

	return nil
}
