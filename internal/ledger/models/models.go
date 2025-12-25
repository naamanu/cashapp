package models

import (
	"cashapp/core"
)

type Transaction struct {
	core.Model
	FailureReason     string             `json:"failure_reason"`
	Direction         core.Direction     `json:"direction"`
	Status            core.Status        `json:"status"`
	Description       string             `json:"description"`
	Ref               string             `json:"ref"`
	From              int                `json:"from"`
	To                int                `json:"to"`
	WalletID          int                `json:"wallet_id"`
	Amount            int64              `json:"amount"`
	Purpose           core.Purpose       `json:"purpose"`
	Privacy           string             `json:"privacy" gorm:"default:'private'"` // public, friends, private
	TransactionEvents []TransactionEvent `json:"transaction_events"`
}

type TransactionEvent struct {
	core.Model
	TransactionID int       `json:"transaction_id"`
	WalletID      int       `json:"wallet_id"`
	Type          core.Type `json:"type"`
	Amount        int64     `json:"amount"`
}

type PaymentRequest struct {
	core.Model
	RequesterID int    `json:"requester_id"`
	PayerID     int    `json:"payer_id"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"` // pending, paid, declined
	Description string `json:"description"`
}
