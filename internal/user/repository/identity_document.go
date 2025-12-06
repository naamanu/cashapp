package repository

import (
	"cashapp/internal/user/models"

	"gorm.io/gorm"
)

type identityDocumentLayer struct {
	db *gorm.DB
}

type IdentityDocumentRepo interface {
	Create(doc *models.IdentityDocument) error
	Update(doc *models.IdentityDocument) error
	FindByID(id int) (*models.IdentityDocument, error)
	FindByUserID(userID int) ([]models.IdentityDocument, error)
}

func newIdentityDocumentLayer(db *gorm.DB) *identityDocumentLayer {
	return &identityDocumentLayer{
		db: db,
	}
}

func (l *identityDocumentLayer) Create(doc *models.IdentityDocument) error {
	return l.db.Create(doc).Error
}

func (l *identityDocumentLayer) Update(doc *models.IdentityDocument) error {
	return l.db.Save(doc).Error
}

func (l *identityDocumentLayer) FindByID(id int) (*models.IdentityDocument, error) {
	var doc models.IdentityDocument
	if err := l.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (l *identityDocumentLayer) FindByUserID(userID int) ([]models.IdentityDocument, error) {
	var docs []models.IdentityDocument
	err := l.db.Where("user_id = ?", userID).Find(&docs).Error
	return docs, err
}
