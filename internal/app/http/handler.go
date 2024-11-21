package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"rest-wallet/internal/app/api"
)

// Handler struct wraps service interface, which interacts with business logic
type Handler struct {
	service api.Service
	logger  *logrus.Logger
}

// New creates new Handler instance and takes api.Service and logger as parameters
func New(service api.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers HTTP routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// API Routes

	// @Router /api/v1/wallet [post]
	r.HandleFunc("/api/v1/wallet", h.UpdateWalletHandler).Methods("POST")

	// @Router /api/v1/wallets/{WALLET_UUID} [get]
	r.HandleFunc("/api/v1/wallets/{WALLET_UUID}", h.GetWalletHandler).Methods("GET")

	// Swagger documentation endpoint
	r.PathPrefix("/docs/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/docs/swagger/index.html", httpSwagger.WrapHandler)
}

// StartServer initializes and starts HTTP server on given port
func (h *Handler) StartServer(port string) {
	router := mux.NewRouter()

	h.RegisterRoutes(router)

	if err := http.ListenAndServe(port, router); err != nil {
		h.logger.Fatalf("Не удалось запустить сервер: %s", err)
	}
}
