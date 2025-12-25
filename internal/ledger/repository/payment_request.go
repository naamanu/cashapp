package repository

import (
	"cashapp/internal/ledger/models"

	"gorm.io/gorm"
)

type paymentRequestLayer struct {
	db *gorm.DB
}

type PaymentRequestRepo interface {
	Create(req *models.PaymentRequest) error
	FindByID(id int) (*models.PaymentRequest, error)
	ListByPayer(payerID int) ([]models.PaymentRequest, error)
	Update(req *models.PaymentRequest) error
}

func newPaymentRequestLayer(db *gorm.DB) *paymentRequestLayer {
	return &paymentRequestLayer{
		db: db,
	}
}

func (l *paymentRequestLayer) Create(req *models.PaymentRequest) error {
	return l.db.Create(req).Error
}

func (l *paymentRequestLayer) FindByID(id int) (*models.PaymentRequest, error) {
	var req models.PaymentRequest
	if err := l.db.First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (l *paymentRequestLayer) ListByPayer(payerID int) ([]models.PaymentRequest, error) {
	var reqs []models.PaymentRequest
	err := l.db.Where("payer_id = ?", payerID).Find(&reqs).Error
	return reqs, err
}

func (l *paymentRequestLayer) Update(req *models.PaymentRequest) error {
	return l.db.Save(req).Error
}
