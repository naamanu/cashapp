package models

import (
	"cashapp/core"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	core.Model
	Tag     string   `json:"tag"`
	Wallets []Wallet `json:"wallets"`
}

type Wallet struct {
	core.Model
	UserID    int   `json:"user_id"`
	User      *User `json:"user,omitempty"`
	IsPrimary bool  `json:"is_primary,omitempty"`
}

func RunSeeds(db *gorm.DB) {
	user := User{
		Tag: "yaw",
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
	}

	if err := db.Model(&Wallet{}).Where("account_id=? AND is_primary=?", user.ID, true).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&wallet)
		} else {
			core.Log.Info("record found or other error", zap.Error(err))
		}
	}

}
