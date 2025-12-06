package core

import "time"

type Type string
type Status string
type Direction string
type Purpose string

var (
	TypeDebit  Type = "debit"
	TypeCredit Type = "credit"

	StatusFailed  Status = "failed"
	StatusPending Status = "pending"
	StatusSuccess Status = "success"

	DirectionIncoming Direction = "incoming"
	DirectionOutgoing Direction = "outgoing"

	PurposeTransfer   Purpose = "transfer"
	PurposeDeposit    Purpose = "deposit"
	PurposeWithdrawal Purpose = "withdrawal"
	PurposeReversal   Purpose = "reversal"
)

type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
