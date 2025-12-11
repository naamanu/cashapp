package repository

import "gorm.io/gorm"

type Repo struct {
	Users   UserRepo
	Wallets WalletRepo
}

func New(db *gorm.DB) Repo {
	return Repo{
		Users:   newUserLayer(db),
		Wallets: newWalletLayer(db),
	}
}
