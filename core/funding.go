package core

type LinkFundingSourceRequest struct {
	UserID          int    `json:"user_id"`
	PaymentMethodID string `json:"payment_method_id"` // "pm_card_visa"
	Type            string `json:"type"`              // card
}

type DepositRequest struct {
	UserID          int   `json:"user_id"`
	Amount          int64 `json:"amount"` // in cents
	FundingSourceID int   `json:"funding_source_id"`
}
