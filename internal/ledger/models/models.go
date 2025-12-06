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
	TransactionEvents []TransactionEvent `json:"transaction_events"`
}

type TransactionEvent struct {
	core.Model
	TransactionID int       `json:"transaction_id"`
	WalletID      int       `json:"wallet_id"`
	Type          core.Type `json:"type"`
	Amount        int64     `json:"amount"`
}
