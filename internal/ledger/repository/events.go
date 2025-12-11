package repository

import (
	"cashapp/core"
	"cashapp/internal/ledger/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type eventLayer struct {
	db *gorm.DB
}

type EventRepo interface {
	GetWalletBalance(id int) (int64, error)
	Save(tx *gorm.DB, data *models.TransactionEvent) error
}

func newEventLayer(db *gorm.DB) *eventLayer {
	return &eventLayer{
		db: db,
	}
}

func (el *eventLayer) GetWalletBalance(id int) (int64, error) {
	var balance int64

	rows, err := el.db.Table("transaction_events").Select("amount, type").Where("wallet_id = ?", id).Rows()
	if err != nil {
		return balance, err
	}
	defer rows.Close()

	for rows.Next() {
		var amount int64
		var event_type string
		if err := rows.Scan(&amount, &event_type); err != nil {
			return 0, fmt.Errorf("error reading amount/type: %v", err)
		}
		if strings.EqualFold(event_type, string(core.TypeDebit)) {
			balance -= amount
		} else {
			balance += amount
		}
	}

	return balance, nil
}

func (el *eventLayer) Save(tx *gorm.DB, data *models.TransactionEvent) error {
	if err := tx.Create(data).Error; err != nil {
		return err
	}
	return nil
}
