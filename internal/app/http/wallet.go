package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"rest-wallet/internal/app/api"
	"rest-wallet/internal/app/models"
)

// UpdateWalletHandler updates a wallet by applying a deposit or withdrawal operation.
// @Summary      Update Wallet
// @Description  Updates a wallet by applying a deposit or withdrawal operation.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        input  body      models.WalletUpdateRequest  true  "Wallet update request"
// @Success      200    {string}  string                      "Wallet updated successfully"
// @Failure      400    {string}  string                      "Invalid input"
// @Failure      404    {string}  string                      "Wallet not found"
// @Failure      500    {string}  string                      "Internal server error"
// @Router /api/v1/wallet [post]
func (h *Handler) UpdateWalletHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("UpdateWalletHandler[http]: Updating wallet")

	// Decode request body into input struct
	var input models.WalletUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	// Check request validity
	if validity := input.Validate(); validity != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %v", validity), http.StatusBadRequest)
		return
	}

	// Attempt to update wallet using service
	if err := h.service.UpdateWallet(input); err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h.logger.Infof("UpdateWalletHandler[http]: Successful wallet update")
}

// GetWalletHandler retrieves the balance of a wallet by its UUID.
// @Summary      Get Wallet Balance
// @Description  Retrieves the balance of a wallet by its UUID.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        WALLET_UUID  path      string  true  "Wallet UUID"
// @Success      200          {object}  models.Wallet  "Wallet balance retrieved successfully"
// @Failure      400          {string}  string         "Invalid walletId"
// @Failure      404          {string}  string         "Wallet not found"
// @Failure      500          {string}  string         "Internal server error"
// @Router /api/v1/wallets/{WALLET_UUID} [get]
func (h *Handler) GetWalletHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("GetWalletHandler[http]: Getting wallet")

	// Get wallet ID from request
	vars := mux.Vars(r)
	walletIDStr := vars["WALLET_UUID"]

	// Check if wallet ID is missing
	if walletIDStr == "" {
		http.Error(w, "Wallet UUID is missing", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("Received WALLET_UUID: %s", walletIDStr)

	// Check UUID validity
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		h.logger.Errorf("Invalid WALLET_UUID: %v", err)
		http.Error(w, "Invalid walletId", http.StatusBadRequest)
		return
	}

	// Attempt to get wallet balance using service
	balance, err := h.service.GetWalletBalance(walletID)
	if err != nil {
		if errors.Is(err, api.ErrWalletNotFound) {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("Responding with balance: %d", balance)

	wallet := models.Wallet{
		ID:      walletID,
		Balance: balance,
	}

	// Respond with wallet balance
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(wallet); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	h.logger.Debugf("GetWalletHandler[http]: Wallet successfully retrieved")
}

// handleServiceError processes service layer errors
func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	errMsg := "Error processing request: %s"
	var statusCode int

	switch {
	case errors.Is(err, api.ErrWalletNotFound):
		errMsg = fmt.Sprintf(errMsg, err)
		statusCode = http.StatusNotFound
	case errors.Is(err, api.ErrInsufficientFunds):
		errMsg = fmt.Sprintf(errMsg, err)
		statusCode = http.StatusBadRequest
	case errors.Is(err, api.ErrExceedMaxBalance):
		errMsg = fmt.Sprintf(errMsg, err)
		statusCode = http.StatusBadRequest
	default:
		errMsg = "Internal server error"
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, errMsg, statusCode)
}
