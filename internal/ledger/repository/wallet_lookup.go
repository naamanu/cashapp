package repository

import (
	"gorm.io/gorm"
)

// Minimal definition for lookup
type walletStub struct {
	ID        int
	UserID    int
	IsPrimary bool
}

type WalletLookupRepo interface {
	GetPrimaryWalletID(userID int) (int, error)
}

type walletLookupLayer struct {
	db *gorm.DB
}

func newWalletLookupLayer(db *gorm.DB) *walletLookupLayer {
	return &walletLookupLayer{db: db}
}

func (l *walletLookupLayer) GetPrimaryWalletID(userID int) (int, error) {
	var w walletStub
	err := l.db.Table("wallets").Select("id").Where("user_id = ? AND is_primary = ?", userID, true).First(&w).Error
	return w.ID, err
}
