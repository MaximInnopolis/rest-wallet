package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"rest-wallet/internal/app/config"
)

var CacheExpiration = 5 * time.Minute

type ClientRedis struct {
	Client *redis.Client
}

func NewRedisClient(conf config.RedisConfig) (*ClientRedis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		PoolSize: conf.PoolSize, // Number of connections in pool
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &ClientRedis{Client: client}, nil
}

// SetWalletBalance sets the wallet balance in Redis
func (r *ClientRedis) SetWalletBalance(ctx context.Context, walletID uuid.UUID, balance int64) error {
	key := fmt.Sprintf("wallet:%s:balance", walletID)
	return r.Client.Set(ctx, key, balance, CacheExpiration).Err()
}

// GetWalletBalance retrieves the wallet balance from Redis
func (r *ClientRedis) GetWalletBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	key := fmt.Sprintf("wallet:%s:balance", walletID)
	balance, err := r.Client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, fmt.Errorf("cache miss")
	}
	return balance, err
}

// Close gracefully closes the Redis client connection.
func (r *ClientRedis) Close() error {
	return r.Client.Close()
}
