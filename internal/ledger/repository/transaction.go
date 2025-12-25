package repository

import (
	"cashapp/internal/ledger/models"

	"gorm.io/gorm"
)

type transactionLayer struct {
	db *gorm.DB
}

type TransactionRepo interface {
	SQLTransaction(f func(tx *gorm.DB) error) error
	Create(tx *gorm.DB, data *models.Transaction) error
	Updates(tx *gorm.DB, transactions ...*models.Transaction) error
	GetFeed(friendIDs []int) ([]models.Transaction, error)
	FindByID(id int) (*models.Transaction, error)
}

func newTransactionLayer(db *gorm.DB) *transactionLayer {
	return &transactionLayer{
		db: db,
	}
}

func (tl *transactionLayer) FindByID(id int) (*models.Transaction, error) {
	var tx models.Transaction
	err := tl.db.First(&tx, id).Error
	return &tx, err
}

func (tl *transactionLayer) GetFeed(friendIDs []int) ([]models.Transaction, error) {
	var txs []models.Transaction
	// Find transactions involving any friend where privacy is not 'private'
	// Simplification: showing all public/friends transactions involving these users.
	err := tl.db.Where("( \"from\" IN ? OR \"to\" IN ? ) AND privacy != 'private'", friendIDs, friendIDs).
		Order("created_at desc").
		Limit(20).
		Find(&txs).Error
	return txs, err
}

func (tl *transactionLayer) SQLTransaction(f func(tx *gorm.DB) error) error {
	return tl.db.Transaction(f)
}

func (tl *transactionLayer) Create(tx *gorm.DB, data *models.Transaction) error {
	if err := tx.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (tl *transactionLayer) Updates(tx *gorm.DB, transactions ...*models.Transaction) error {
	for _, trans := range transactions {
		if err := tx.Updates(trans).Error; err != nil {
			return err
		}
	}
	return nil
}
