package processor

import (
	"cashapp/core"
	"cashapp/internal/ledger/models"
	"cashapp/internal/ledger/repository"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Processor struct {
	Repo repository.Repo
}

func New(r repository.Repo) Processor {
	return Processor{
		Repo: r,
	}
}

func (p *Processor) ProcessTransaction(fromTrans models.Transaction) error {
	switch fromTrans.Purpose {
	case core.PurposeTransfer:
		f, t, err := p.MoveMoneyBetweenWallets(fromTrans)
		if err != nil {
			if err := p.FailureCallback(f, t, err); err != nil {
				return fmt.Errorf("failed to complete transaction. %v", err)
			}
			return fmt.Errorf("money transfer failed. %v", err)
		}
		if err := p.SuccessCallback(f, t); err != nil {
			return fmt.Errorf("failed to complete transaction. %v", err)
		}

	case core.PurposeWithdrawal: // Fixed duplicate case
		if err := p.WithdrawMoneyFromWallet(fromTrans); err != nil {
			return fmt.Errorf("money withdrawal failed. %v", err)
		}
	case core.PurposeDeposit:
		if err := p.DepositMoneyIntoWallet(fromTrans); err != nil {
			return fmt.Errorf("money deposit failed. %v", err)
		}
	default:
		core.Log.Warn("no handler for purpose", zap.Any("purpose", fromTrans.Purpose))
	}
	return nil
}

func (p *Processor) SuccessCallback(fromTrans, toTrans *models.Transaction) error {
	fromTrans.Status = core.StatusSuccess
	toTrans.Status = core.StatusSuccess

	return p.Repo.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		return p.Repo.Transactions.Updates(tx, fromTrans, toTrans)
	})
}

func (p *Processor) FailureCallback(fromTrans, toTrans *models.Transaction, err error) error {
	fromTrans.Status = core.StatusFailed
	toTrans.Status = core.StatusFailed
	fromTrans.FailureReason = err.Error()
	toTrans.FailureReason = err.Error()

	return p.Repo.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		return p.Repo.Transactions.Updates(tx, fromTrans, toTrans)
	})
}
