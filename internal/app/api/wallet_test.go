package api

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"rest-wallet/internal/app/models"
	mockrepo "rest-wallet/internal/app/repository/mocks"
	"rest-wallet/internal/app/repository/postgresql"
)

func TestWalletService_UpdateWallet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mockrepo.NewMockWalletRepo(ctrl)
	logger := logrus.New()

	walletService := NewWalletService(mockWalletRepo, logger)

	tests := []struct {
		name           string
		request        models.WalletUpdateRequest
		mockGetBalance func(walletID uuid.UUID) (int64, error)
		mockUpdate     func(walletID uuid.UUID, amount int64) error
		expectedError  error
	}{
		{
			name: "wallet not found",
			request: models.WalletUpdateRequest{
				WalletID:      uuid.New(),
				Amount:        100,
				OperationType: models.OperationDeposit,
			},
			mockGetBalance: func(walletID uuid.UUID) (int64, error) {
				return 0, postgresql.ErrNotFound
			},
			mockUpdate:    nil,
			expectedError: ErrWalletNotFound,
		},
		{
			name: "deposit exceeds max balance",
			request: models.WalletUpdateRequest{
				WalletID:      uuid.New(),
				Amount:        10000001,
				OperationType: models.OperationDeposit,
			},
			mockGetBalance: func(walletID uuid.UUID) (int64, error) {
				return 0, nil
			},
			mockUpdate:    nil,
			expectedError: ErrExceedMaxBalance,
		},
		{
			name: "insufficient funds for withdrawal",
			request: models.WalletUpdateRequest{
				WalletID:      uuid.New(),
				Amount:        500,
				OperationType: models.OperationWithdraw,
			},
			mockGetBalance: func(walletID uuid.UUID) (int64, error) {
				return 100, nil // Balance is less than the withdrawal amount
			},
			mockUpdate:    nil,
			expectedError: ErrInsufficientFunds,
		},
		{
			name: "successful deposit",
			request: models.WalletUpdateRequest{
				WalletID:      uuid.New(),
				Amount:        100,
				OperationType: models.OperationDeposit,
			},
			mockGetBalance: func(walletID uuid.UUID) (int64, error) {
				return 500, nil // Initial balance is 500
			},
			mockUpdate: func(walletID uuid.UUID, amount int64) error {
				return nil // Simulate successful update
			},
			expectedError: nil,
		},
		{
			name: "successful withdrawal",
			request: models.WalletUpdateRequest{
				WalletID:      uuid.New(),
				Amount:        100,
				OperationType: models.OperationWithdraw,
			},
			mockGetBalance: func(walletID uuid.UUID) (int64, error) {
				return 500, nil // Initial balance is 500
			},
			mockUpdate: func(walletID uuid.UUID, amount int64) error {
				return nil // Simulate successful update
			},
			expectedError: nil,
		},
	}

	// Running test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			mockWalletRepo.
				EXPECT().GetWalletBalance(tt.request.WalletID).DoAndReturn(tt.mockGetBalance).Times(1)
			if tt.mockUpdate != nil {
				mockWalletRepo.
					EXPECT().Update(tt.request.WalletID, gomock.Any()).DoAndReturn(tt.mockUpdate).Times(1)
			}

			// act
			err := walletService.UpdateWallet(tt.request)

			// assert
			if err != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error: %v, but got nil", tt.expectedError)
			}
		})
	}

}

func TestWalletService_GetWalletBalance(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletRepo := mockrepo.NewMockWalletRepo(ctrl)
	logger := logrus.New()

	walletService := NewWalletService(mockWalletRepo, logger)

	tests := []struct {
		name            string
		walletID        uuid.UUID
		expectedBalance int64
		expectedError   error
		setupMocks      func(wallerID uuid.UUID)
	}{
		{
			name:     "success",
			walletID: uuid.New(),
			setupMocks: func(walletID uuid.UUID) {
				mockWalletRepo.
					EXPECT().GetWalletBalance(walletID).Return(int64(100), nil)
			},
			expectedBalance: 100,
			expectedError:   nil,
		},
		{
			name:     "wallet not found",
			walletID: uuid.New(),
			setupMocks: func(walletID uuid.UUID) {
				mockWalletRepo.
					EXPECT().GetWalletBalance(walletID).Return(int64(0), postgresql.ErrNotFound)
			},
			expectedBalance: 0,
			expectedError:   ErrWalletNotFound,
		},
	}

	// Running test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			tt.setupMocks(tt.walletID)

			// act
			balance, err := walletService.GetWalletBalance(tt.walletID)

			// assert
			assert.Equal(t, tt.expectedBalance, balance)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError, "Expected error does not match actual error")
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
			}
		})
	}
}
