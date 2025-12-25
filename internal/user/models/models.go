package models

import (
	"cashapp/core"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusVerified KYCStatus = "verified"
	KYCStatusRejected KYCStatus = "rejected"
)

type User struct {
	core.Model
	Tag            string    `json:"tag"`
	Wallets        []Wallet  `json:"wallets"`
	KYCLevel       int       `json:"kyc_level"` // 0: Unverified, 1: Basic, 2: Full
	KYCStatus      KYCStatus `json:"kyc_status" gorm:"default:'pending'"`
	RiskScore      int       `json:"risk_score"`
	DefaultPrivacy string    `json:"default_privacy" gorm:"default:'public'"` // public, friends, private
}

type Friendship struct {
	core.Model
	UserID   int    `json:"user_id"`
	FriendID int    `json:"friend_id"`
	Status   string `json:"status"` // pending, accepted
}

type Wallet struct {
	core.Model
	UserID    int    `json:"user_id"`
	User      *User  `json:"user,omitempty"`
	IsPrimary bool   `json:"is_primary,omitempty"`
	Balance   int64  `json:"balance" gorm:"default:0"` // in cents
	Currency  string `json:"currency" gorm:"default:'USD'"`
}

type FundingSource struct {
	core.Model
	UserID     int    `json:"user_id"`
	Type       string `json:"type"`        // card, bank_account
	ProviderID string `json:"provider_id"` // stripe_pm_123
	Last4      string `json:"last4"`
	Brand      string `json:"brand"` // visa, mastercard
}

func RunSeeds(db *gorm.DB) {
	user := User{
		Tag:       "yaw",
		KYCLevel:  1,
		KYCStatus: KYCStatusVerified,
	}

	if err := db.Model(&User{}).Where("tag=?", user.Tag).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&user)
		} else {
			core.Log.Info("record found or other error", zap.Error(err))
		}
	}

	wallet := Wallet{
		UserID:    user.ID,
		IsPrimary: true,
		Balance:   100000, // $1000.00 initial balance for seed user
		Currency:  "USD",
	}

	if err := db.Model(&Wallet{}).Where("account_id=? AND is_primary=?", user.ID, true).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&wallet)
		} else {
			core.Log.Info("record found or other error", zap.Error(err))
		}
	}

}

type IdentityDocument struct {
	core.Model
	UserID int    `json:"user_id"`
	User   *User  `json:"user,omitempty"`
	Type   string `json:"type"`   // passport, drivers_license
	URL    string `json:"url"`    // S3 URL
	Status string `json:"status"` // pending, verified, rejected
}
