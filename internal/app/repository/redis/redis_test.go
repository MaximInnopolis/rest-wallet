package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetWalletBalance(t *testing.T) {
	// arrange
	client, mock := redismock.NewClientMock()
	redisClient := &ClientRedis{Client: client}

	walletID := uuid.New()
	balance := int64(1000)
	mock.ExpectSet(fmt.Sprintf("wallet:%s:balance", walletID), balance, CacheExpiration).SetVal("OK")

	// act
	err := redisClient.SetWalletBalance(context.Background(), walletID, balance)

	// assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWalletBalance(t *testing.T) {
	// arrange
	client, mock := redismock.NewClientMock()
	redisClient := &ClientRedis{Client: client}

	walletID := uuid.New()
	expectedBalance := int64(1000)
	mock.ExpectGet(fmt.Sprintf("wallet:%s:balance", walletID)).SetVal(strconv.FormatInt(expectedBalance, 10))

	// act
	balance, err := redisClient.GetWalletBalance(context.Background(), walletID)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWalletBalanceCacheMiss(t *testing.T) {
	// arrange
	client, mock := redismock.NewClientMock()
	redisClient := &ClientRedis{Client: client}

	walletID := uuid.New()
	mock.ExpectGet(fmt.Sprintf("wallet:%s:balance", walletID)).RedisNil()

	// act
	balance, err := redisClient.GetWalletBalance(context.Background(), walletID)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "cache miss", err.Error())
	assert.Equal(t, int64(0), balance)
	assert.NoError(t, mock.ExpectationsWereMet())
}
