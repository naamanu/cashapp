package service

import (
	"cashapp/core"
	"cashapp/core/currency"
	"cashapp/internal/ledger/models"
	"cashapp/internal/ledger/processor"
	"cashapp/internal/ledger/repository"

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
