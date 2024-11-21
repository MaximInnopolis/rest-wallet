package config

import (
	"fmt"
	"os"
	"strconv"
)

var defaultHttpPort = ":8080"

// RedisConfig struct holds configuration values for Redis address and pool size
type RedisConfig struct {
	Addr     string
	PoolSize int
}

// Config struct holds configuration values for database url, http port and RedisConfig
type Config struct {
	DbUrl    string
	HttpPort string
	RedisConfig
}

// New creates new Config instance by reading environment variables
// It checks if required DATABASE_URL is set; if not, it returns error
// If HTTP_PORT is not set, it defaults to ":8080".
func New() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL не задан")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = defaultHttpPort
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR не задан")
	}

	redisPoolSizeStr := os.Getenv("REDIS_POOL_SIZE")
	if redisPoolSizeStr == "" {
		redisPoolSizeStr = "10"
	}

	redisPoolSize, err := strconv.Atoi(redisPoolSizeStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при конвертации REDIS_POOL_SIZE в число: %v", err)
	}

	return &Config{
		DbUrl:    dbURL,
		HttpPort: httpPort,
		RedisConfig: RedisConfig{
			Addr:     redisAddr,
			PoolSize: redisPoolSize,
		},
	}, nil
}
