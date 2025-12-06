package repository

import (
	"cashapp/internal/user/models"

	"gorm.io/gorm"
)

type walletLayer struct {
	db *gorm.DB
}

type WalletRepo interface {
	Create(userId int) (*models.Wallet, error)
	FindPrimaryWallet(userId int) (*models.Wallet, error)
}

func newWalletLayer(db *gorm.DB) *walletLayer {
	return &walletLayer{
		db: db,
	}
}

func (wl *walletLayer) Create(userId int) (*models.Wallet, error) {
	wallet := models.Wallet{
		UserID:    userId,
		IsPrimary: true,
	}

	if err := wl.db.Create(&wallet).Error; err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (wl *walletLayer) FindPrimaryWallet(userId int) (*models.Wallet, error) {
	wallet := models.Wallet{
		UserID:    userId,
		IsPrimary: true,
	}
	if err := wl.db.First(&wallet).Error; err != nil {
		return nil, err
	}

	return &wallet, nil
}
