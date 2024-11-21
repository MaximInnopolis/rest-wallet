package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	// arrange
	tests := []struct {
		name          string
		request       WalletUpdateRequest
		expectedError error
	}{
		{
			name: "Valid request",
			request: WalletUpdateRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        500,
			},
			expectedError: nil,
		},
		{
			name: "Invalid wallet ID",
			request: WalletUpdateRequest{
				WalletID:      uuid.Nil,
				OperationType: OperationDeposit,
				Amount:        500,
			},
			expectedError: ErrInvalidWalletID,
		},
		{
			name: "Invalid amount (zero)",
			request: WalletUpdateRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        0,
			},
			expectedError: ErrInvalidAmount,
		},
		{
			name: "Invalid amount (negative)",
			request: WalletUpdateRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        -10,
			},
			expectedError: ErrInvalidAmount,
		},
		{
			name: "Amount exceeds maximum",
			request: WalletUpdateRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        2000000, // exceeds maxAmount
			},
			expectedError: ErrExceedMaxAmount,
		},
		{
			name: "Invalid operation type",
			request: WalletUpdateRequest{
				WalletID:      uuid.New(),
				OperationType: "INVALID",
				Amount:        100,
			},
			expectedError: ErrInvalidOperationType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// act
			err := tt.request.Validate()

			// assert
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
