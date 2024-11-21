package api

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rest-wallet/internal/app/models"
	"rest-wallet/internal/app/repository"
	"rest-wallet/internal/app/repository/postgresql"
)

const maxBalance = 10000000

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrExceedMaxBalance  = errors.New("exceed maximum balance")
)

type WalletService struct {
	repo   repository.WalletRepo
	logger *logrus.Logger
}

func NewWalletService(repo repository.WalletRepo, logger *logrus.Logger) *WalletService {
	return &WalletService{
		repo:   repo,
		logger: logger,
	}
}

func (w *WalletService) UpdateWallet(request models.WalletUpdateRequest) error {
	w.logger.Debugf("UpdateWallet[service]: Updating wallet with id: %s", request.WalletID)

	// Ensure that wallet exists
	balance, err := w.repo.GetWalletBalance(request.WalletID)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotFound) {
			w.logger.Errorf("UpdateWallet[service]: Wallet with id not found: %s", request.WalletID)
			return ErrWalletNotFound
		}
		w.logger.Errorf("UpdateWallet[service]: Error retrieving balance for wallet with id: %s: %s", request.WalletID, err)
		return err
	}

	switch request.OperationType {
	case models.OperationDeposit:
		w.logger.Debugf("UpdateWallet[service]: Depositing %d to wallet with id: %s", request.Amount, request.WalletID)

		// Check for balance boundaries
		updatedBalance := balance + request.Amount
		if updatedBalance > maxBalance {
			w.logger.Errorf("UpdateWallet[service]: Wallet balance exceeds maximum: %d", updatedBalance)
			return ErrExceedMaxBalance
		}

		// Deposit amount to wallet
		err = w.repo.Update(request.WalletID, updatedBalance)
		if err != nil {
			w.logger.Errorf("UpdateWallet[service]: Error depositing to wallet with id: %s: %s", request.WalletID, err)
			return err
		}

	case models.OperationWithdraw:
		w.logger.Debugf("UpdateWallet[service]: Withdrawing %d from wallet with id: %s", request.Amount, request.WalletID)

		// Ensure that wallet has sufficient funds
		if balance < request.Amount {
			w.logger.Errorf("UpdateWallet[service]: Insufficient funds for wallet with id: %s. Balance: %d", request.WalletID, balance)
			return ErrInsufficientFunds
		}

		// Withdraw amount from wallet
		updatedBalance := balance - request.Amount
		err = w.repo.Update(request.WalletID, updatedBalance)
		if err != nil {
			w.logger.Errorf("UpdateWallet[service]: Error withdrawing from wallet with id: %s: %s", request.WalletID, err)
			return err
		}

	default:
		w.logger.Errorf("UpdateWallet[service]: Invalid operation type: %s", request.OperationType)
		return models.ErrInvalidOperationType
	}

	w.logger.Infof("UpdateWallet[service]: Wallet with id: %s updated successfully", request.WalletID)
	return nil
}

func (w *WalletService) GetWalletBalance(walletID uuid.UUID) (int64, error) {
	w.logger.Debugf("GetWalletBalance[service]: Getting balance for wallet with id: %s", walletID)

	balance, err := w.repo.GetWalletBalance(walletID)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotFound) {
			w.logger.Errorf("GetWalletBalance[service]: Wallet with id not found: %s: %s", walletID, err)
			return 0, ErrWalletNotFound
		}

		w.logger.Errorf("GetWalletBalance[service]: Error getting balance for wallet with id: %s: %s", walletID, err)
		return 0, err
	}

	w.logger.Infof("GetWalletBalance[service]: Balance for wallet with id %s: %d", walletID, balance)
	return balance, nil
}
