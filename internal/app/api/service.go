package api

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rest-wallet/internal/app/models"
	"rest-wallet/internal/app/repository"
)

type Wallet interface {
	UpdateWallet(request models.WalletUpdateRequest) error
	GetWalletBalance(walletID uuid.UUID) (int64, error)
}

type Service struct {
	Wallet
}

func New(repo *repository.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Wallet: NewWalletService(repo.WalletRepo, logger),
	}
}
