package processor

import (
	"cashapp/core"
	"cashapp/core/currency"
	"cashapp/internal/ledger/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (p *Processor) MoveMoneyBetweenWallets(fromTrans models.Transaction) (*models.Transaction, *models.Transaction, error) {

	originWalletID, err := p.Repo.WalletLookup.GetPrimaryWalletID(fromTrans.From)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find primary wallet for origin. %v", err)
	}

	destinationWalletID, err := p.Repo.WalletLookup.GetPrimaryWalletID(fromTrans.To)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find primary wallet for destination. %v", err)
	}

	balance, err := p.Repo.TransactionEvents.GetWalletBalance(originWalletID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load balance. %v", err)
	}

	if balance < fromTrans.Amount {
		return nil, nil, errors.New("insufficient balance")
	}

	toTrans := models.Transaction{
		From:        fromTrans.From,
		To:          fromTrans.To,
		Ref:         fromTrans.Ref,
		Amount:      currency.ConvertCedisToPessewas(fromTrans.Amount),
		Description: fromTrans.Description,
		Direction:   core.DirectionIncoming,
		Status:      core.StatusPending,
		Purpose:     core.PurposeTransfer,
	}

	err = p.Repo.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		return p.Repo.Transactions.Create(tx, &toTrans)
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create destination transaction. %v", err)
	}

	err = p.Repo.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		debit := models.TransactionEvent{
			TransactionID: fromTrans.ID,
			WalletID:      originWalletID,
			Amount:        fromTrans.Amount,
			Type:          core.TypeDebit,
		}

		if err := p.Repo.TransactionEvents.Save(tx, &debit); err != nil {
			return err
		}

		credit := models.TransactionEvent{
			TransactionID: toTrans.ID,
			WalletID:      destinationWalletID,
			Amount:        toTrans.Amount,
			Type:          core.TypeCredit,
		}

		if err := p.Repo.TransactionEvents.Save(tx, &credit); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("money movement failed. err=%v", err)
	}

	fromTrans.WalletID = originWalletID
	toTrans.WalletID = destinationWalletID

	return &fromTrans, &toTrans, nil
}

func (p *Processor) DepositMoneyIntoWallet(fromTrans models.Transaction) error {
	return nil
}

func (p *Processor) WithdrawMoneyFromWallet(fromTrans models.Transaction) error {
	return nil
}
