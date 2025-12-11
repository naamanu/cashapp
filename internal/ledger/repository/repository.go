package repository

import "gorm.io/gorm"

type Repo struct {
	Transactions      TransactionRepo
	TransactionEvents EventRepo
	WalletLookup      WalletLookupRepo
}

func New(db *gorm.DB) Repo {
	return Repo{
		Transactions:      newTransactionLayer(db),
		TransactionEvents: newEventLayer(db),
		WalletLookup:      newWalletLookupLayer(db),
	}
}
