package repository

import "gorm.io/gorm"

type Repo struct {
	Users             UserRepo
	Wallets           WalletRepo
	IdentityDocuments IdentityDocumentRepo
	FundingSources    FundingSourceRepo
}

func New(db *gorm.DB) Repo {
	return Repo{
		Users:             newUserLayer(db),
		Wallets:           newWalletLayer(db),
		IdentityDocuments: newIdentityDocumentLayer(db),
		FundingSources:    newFundingSourceLayer(db),
	}
}
