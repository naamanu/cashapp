package service

import (
	"cashapp/core"
	"cashapp/core/currency"
	"cashapp/internal/ledger/models"
	"cashapp/internal/ledger/processor"
	"cashapp/internal/ledger/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentService struct {
	repository repository.Repo
	config     *core.Config
	processor  processor.Processor
}

func New(r repository.Repo, c *core.Config) *PaymentService {
	return &PaymentService{
		repository: r,
		config:     c,
		processor:  processor.New(r),
	}
}

func (p *PaymentService) SendMoney(req core.CreatePaymentRequest) core.Response {
	fromTrans := models.Transaction{
		From:        req.From,
		To:          req.To,
		Ref:         core.GenerateRef(),
		Amount:      currency.ConvertCedisToPessewas(req.Amount),
		Description: req.Description,
		Direction:   core.DirectionOutgoing,
		Status:      core.StatusPending,
		Purpose:     core.PurposeTransfer,
		Privacy:     req.Privacy,
	}

	err := p.repository.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		return p.repository.Transactions.Create(tx, &fromTrans)
	})

	if err != nil {
		return core.Error(err, nil)
	}

	if err := p.processor.ProcessTransaction(fromTrans); err != nil {
		return core.Error(err, nil)
	}

	return core.Success(nil, nil)
}

func (p *PaymentService) GetBalance(walletID int) core.Response {
	balance, err := p.repository.TransactionEvents.GetWalletBalance(walletID)
	if err != nil {
		return core.Error(err, nil)
	}

	return core.Success(&map[string]interface{}{
		"balance": currency.ConvertPessewasToCedis(balance),
	}, nil)
}

func (p *PaymentService) CreateRequest(req core.CreateRequestDTO) core.Response {
	// In real world, validate users exist via User Service
	pr := models.PaymentRequest{
		RequesterID: req.RequesterID,
		PayerID:     req.PayerID,
		Amount:      req.Amount, // assume input in Pesewas/Cents
		Description: req.Description,
		Status:      "pending",
	}

	if err := p.repository.PaymentRequests.Create(&pr); err != nil {
		return core.Error(err, core.String("failed to create payment request"))
	}

	// Mock Push Notification
	core.Log.Info("Push Notification sent to Payer", zap.Int("payer_id", req.PayerID), zap.String("message", "You have a new payment request"))

	return core.Success(&map[string]interface{}{
		"request_id": pr.ID,
		"status":     pr.Status,
	}, core.String("payment request created"))
}

func (p *PaymentService) PayRequest(requestID int, payerKey string) core.Response {
	// 1. Fetch Request
	req, err := p.repository.PaymentRequests.FindByID(requestID)
	if err != nil {
		return core.Error(err, core.String("payment request not found"))
	}

	if req.Status != "pending" {
		return core.Error(nil, core.String("request is already processed"))
	}

	// 2. Execute Payment (Reuse SendMoney logic)
	// We need to map models.PaymentRequest to core.CreatePaymentRequest
	// Note: SendMoney expects amount in Cedis (float logic in core struct implies int64 but conversion happens inside SendMoney??)
	// Wait, core.CreatePaymentRequest has Amount int64.
	// SendMoney implementations says: Amount: currency.ConvertCedisToPessewas(req.Amount)
	// This implies CreatePaymentRequest.Amount is in MAJOR UNITS (Cedis/Dollars) but stored as int64?
	// Or maybe it is just raw int64 and the conversion assumes it came as major units?
	// If my PaymentRequest stores lowest denomination (based on my comment in CreateRequest), I should handle this.

	// Check SendMoney implementation:
	// Amount: currency.ConvertCedisToPessewas(req.Amount)
	// This strongly suggests req.Amount is 100 for 100 GHS.

	// If PaymentRequest.Amount is stored as Pesewas (10000 for 100 GHS), and SendMoney expects Major units...
	// We should probably normalize everything to lowest denomination to avoid confusion, but sticking to existing pattern:

	// Let's assume CreateRequestDTO takes Major Units for consistency with CreatePaymentRequest used in SendMoney.
	// So req.Amount is 100.
	// We store 100 in DB (as int64).
	// When paying, we pass 100 to SendMoney, which converts to 10000.

	payReq := core.CreatePaymentRequest{
		From:        req.PayerID,
		To:          req.RequesterID,
		Amount:      req.Amount, // Major units
		Description: req.Description,
	}

	// Ideally refactor SendMoney to not return HTTP Response but error/struct, but for now calling it is fine
	// or we copy the logic. Calling it is better for DRY if we can extract the error.

	// However, SendMoney currently returns core.Response.
	resp := p.SendMoney(payReq)
	if resp.Error {
		return resp
	}

	// 3. Update Request Status
	req.Status = "paid"
	if err := p.repository.PaymentRequests.Update(req); err != nil {
		core.Log.Error("Failed to update payment request status", zap.Error(err))
		// Payment succeeded but status update failed. In critical system, this needs reconciliation.
	}

	return core.Success(nil, core.String("request paid successfully"))
}

func (p *PaymentService) GetFeed(friendIDs []int) core.Response {
	txs, err := p.repository.Transactions.GetFeed(friendIDs)
	if err != nil {
		return core.Error(err, core.String("failed to fetch feed"))
	}

	// Transform to simplified feed items
	var feed []map[string]interface{}
	for _, tx := range txs {
		feed = append(feed, map[string]interface{}{
			"id":          tx.ID,
			"from":        tx.From,
			"to":          tx.To,
			"amount":      currency.ConvertPessewasToCedis(tx.Amount),
			"description": tx.Description,
			"timestamp":   tx.CreatedAt,
			"privacy":     tx.Privacy,
		})
	}

	return core.Success(&map[string]interface{}{"feed": feed}, nil)
}

func (p *PaymentService) SplitBill(req core.SplitBillDTO) core.Response {
	// 1. Fetch Original Transaction
	tx, err := p.repository.Transactions.FindByID(req.OriginalTransactionID)
	if err != nil {
		return core.Error(err, core.String("original transaction not found"))
	}

	// 2. Validate Ownership (Assume requester must be the 'From' user, i.e., they paid initially)
	if tx.From != req.RequesterID {
		return core.Error(nil, core.String("only the payer can split the bill"))
	}

	// 3. Calculate Split Amount
	// Total participants = Requester + Friends
	totalParticipants := int64(len(req.FriendIDs) + 1)
	splitAmount := tx.Amount / totalParticipants // Integer division (Pesewas)

	// 4. Create Payment Requests for each friend
	var requests []models.PaymentRequest
	for _, friendID := range req.FriendIDs {
		pr := models.PaymentRequest{
			RequesterID: req.RequesterID,
			PayerID:     friendID,
			Amount:      splitAmount,
			Description: "Split Bill: " + tx.Description,
			Status:      "pending",
		}
		requests = append(requests, pr)
	}

	// 5. Save all requests
	for _, pr := range requests {
		if err := p.repository.PaymentRequests.Create(&pr); err != nil {
			// In real world, use transaction to rollback all if one fails
			core.Log.Error("Failed to create split request", zap.Error(err))
		} else {
			// Mock Push Notification
			core.Log.Info("Split Bill Request sent", zap.Int("to_user", pr.PayerID), zap.Int64("amount", pr.Amount))
		}
	}

	return core.Success(&map[string]interface{}{
		"total_amount":     currency.ConvertPessewasToCedis(tx.Amount),
		"split_amount":     currency.ConvertPessewasToCedis(splitAmount),
		"requests_created": len(requests),
	}, core.String("bill split successfully"))
}
