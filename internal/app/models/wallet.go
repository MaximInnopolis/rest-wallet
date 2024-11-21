package models

import "github.com/google/uuid"

// Wallet represents wallet model
type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance int64     `json:"balance"`
}
