package core

type Pagination struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	NextPage     int   `json:"next_page,omitempty"`
	PreviousPage int   `json:"previous_page,omitempty"`
	Count        int64 `json:"count"`
}

type Meta struct {
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Message    string      `json:"message"`
}

type Response struct {
	Error bool `json:"error"`
	Code  int  `json:"code"`
	Meta  Meta `json:"meta"`
}

type CreateUserRequest struct {
	Tag string `json:"tag"`
}

type CreateFriendshipRequest struct {
	UserID   int `json:"user_id"`
	FriendID int `json:"friend_id"`
}

type CreatePaymentRequest struct {
	From        int    `json:"from"`
	To          int    `json:"to"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Privacy     string `json:"privacy,omitempty"` // public, friends, private
}

type CreateRequestDTO struct {
	RequesterID int    `json:"requester_id"`
	PayerID     int    `json:"payer_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type SplitBillDTO struct {
	OriginalTransactionID int   `json:"original_transaction_id"`
	RequesterID           int   `json:"requester_id"`
	FriendIDs             []int `json:"friend_ids"`
}
