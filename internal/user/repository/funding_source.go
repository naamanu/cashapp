package repository

import (
	"cashapp/internal/user/models"

	"gorm.io/gorm"
)

type fundingSourceLayer struct {
	db *gorm.DB
}

type FundingSourceRepo interface {
	Create(fs *models.FundingSource) error
	ListByUserID(userID int) ([]models.FundingSource, error)
	FindByID(id int) (*models.FundingSource, error)
}

func newFundingSourceLayer(db *gorm.DB) *fundingSourceLayer {
	return &fundingSourceLayer{
		db: db,
	}
}

func (l *fundingSourceLayer) Create(fs *models.FundingSource) error {
	return l.db.Create(fs).Error
}

func (l *fundingSourceLayer) ListByUserID(userID int) ([]models.FundingSource, error) {
	var sources []models.FundingSource
	err := l.db.Where("user_id = ?", userID).Find(&sources).Error
	return sources, err
}

func (l *fundingSourceLayer) FindByID(id int) (*models.FundingSource, error) {
	var fs models.FundingSource
	if err := l.db.First(&fs, id).Error; err != nil {
		return nil, err
	}
	return &fs, nil
}
