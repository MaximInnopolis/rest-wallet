package main

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	_ "rest-wallet/docs"
	"rest-wallet/internal/app/api"
	"rest-wallet/internal/app/config"
	httpHandler "rest-wallet/internal/app/http"
	"rest-wallet/internal/app/repository"
	"rest-wallet/internal/app/repository/database"
	"rest-wallet/internal/app/repository/redis"
)

// @title Wallet API
// @version 1.0
// @description REST API service for wallet
// @host localhost:8080
// @basePath /
// @schemes http

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	// Create config
	cfg, err := config.New()
	if err != nil {
		log.Errorf("Ошибка при чтении конфига: %v", err)
		os.Exit(1)
	}

	// Create a new connection pool to database
	pool, err := database.NewPool(cfg.DbUrl)
	if err != nil {
		log.Errorf("Ошибка при создании соединения к базе данных: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	redisClient, err := redis.NewRedisClient(cfg.RedisConfig)
	if err != nil {
		log.Errorf("Ошибка при создании соединения к Redis: %v", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Create a new Database with connection pool
	db := database.NewDatabase(pool)

	// Create a new repo with Database and logger
	repo := repository.New(*db, redisClient, log)

	// Create a new service
	walletService := api.New(repo, log)

	// Create Http handler
	handler := httpHandler.New(*walletService, log)

	// Init Router
	r := mux.NewRouter()

	handler.RegisterRoutes(r)

	// Start server
	handler.StartServer(cfg.HttpPort)
}
