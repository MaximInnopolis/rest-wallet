package repository

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rest-wallet/internal/app/repository/database"
	"rest-wallet/internal/app/repository/postgresql"
	"rest-wallet/internal/app/repository/redis"
)

type WalletRepo interface {
	Update(walletID uuid.UUID, amount int64) error
	GetWalletBalance(walletID uuid.UUID) (int64, error)
}

type Repository struct {
	WalletRepo
}

func New(db database.Database, cache *redis.ClientRedis, logger *logrus.Logger) *Repository {
	return &Repository{
		WalletRepo: postgresql.NewWalletPostgres(db, cache, logger),
	}
}
