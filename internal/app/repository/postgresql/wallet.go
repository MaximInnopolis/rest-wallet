package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"rest-wallet/internal/app/repository/database"
	"rest-wallet/internal/app/repository/redis"
)

var ErrNotFound = errors.New("wallet not found")

type WalletPostgres struct {
	db     database.Database
	cache  *redis.ClientRedis
	logger *logrus.Logger
}

func NewWalletPostgres(db database.Database, cache *redis.ClientRedis, logger *logrus.Logger) *WalletPostgres {
	return &WalletPostgres{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (w *WalletPostgres) Update(walletID uuid.UUID, amount int64) error {
	w.logger.Debugf("Update[repo]: Updating wallet with id: %s, amount: %d", walletID, amount)

	query := `SELECT 1 FROM wallets WHERE id = $1 FOR UPDATE` // Use to lock row for update and avoid race condition

	ctx := context.Background()

	// Begin transaction
	tx, err := w.db.GetPool().Begin(ctx)
	if err != nil {
		w.logger.Errorf("Update[repo]: Error beginning transaction: %s", err)
		return err
	}
	defer tx.Rollback(ctx) // Rollback transaction if function returns error

	// Lock row for update
	_, err = tx.Exec(ctx, query, walletID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.logger.Errorf("Update[repo]: Wallet with id not found: %s", walletID)
			return ErrNotFound
		}

		w.logger.Errorf("Update[repo]: Error locking wallet with id: %s: %s", walletID, err)
		return err
	}

	updateQuery := `UPDATE wallets SET balance = $1 WHERE id = $2`

	// Update wallet balance
	_, err = tx.Exec(ctx, updateQuery, amount, walletID)
	if err != nil {
		w.logger.Errorf("Update[repo]: Error updating wallet with id: %s: %s", walletID, err)
		return err
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		w.logger.Errorf("Update[repo]: Error committing transaction: %s", err)
		return err
	}

	// Update Redis cache
	err = w.cache.SetWalletBalance(ctx, walletID, amount)
	if err != nil {
		w.logger.Warnf("Update[repo]: Failed to update Redis cache for wallet %s: %s", walletID, err)
	}

	w.logger.Debugf("Update[repo]: Wallet with id %s updated. New balance: %d", walletID, amount)
	return nil
}

func (w *WalletPostgres) GetWalletBalance(walletID uuid.UUID) (int64, error) {
	w.logger.Debugf("GetWalletBalance[repo]: Getting balance for wallet with id: %s", walletID)

	ctx := context.Background()

	// Check Redis cache
	cachedBalance, err := w.cache.GetWalletBalance(ctx, walletID)
	if err == nil {
		w.logger.Debugf("GetWalletBalance[repo]: Cache hit for wallet %s: %d", walletID, cachedBalance)
		return cachedBalance, nil
	}
	w.logger.Warnf("GetWalletBalance[repo]: Cache miss for wallet %s", walletID)

	// if cache miss, get balance from database
	query := `SELECT balance FROM wallets WHERE id = $1`
	var balance int64

	err = w.db.GetPool().QueryRow(ctx, query, walletID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.logger.Errorf("GetWalletBalance[repo]: Wallet with id not found: %s: %s", walletID, err)
			return 0, ErrNotFound
		}

		w.logger.Errorf("GetWalletBalance[repo]: Error getting balance for wallet with id: %s: %s", walletID, err)
		return 0, err
	}

	// Update Redis cache
	err = w.cache.SetWalletBalance(ctx, walletID, balance)
	if err != nil {
		w.logger.Warnf("GetWalletBalance[repo]: Failed to update Redis cache for wallet %s: %s", walletID, err)
	}

	w.logger.Debugf("GetWalletBalance[repo]: Balance for wallet with id %s: %d", walletID, balance)
	return balance, nil
}
