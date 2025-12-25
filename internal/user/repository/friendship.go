package repository

import (
	"cashapp/internal/user/models"

	"gorm.io/gorm"
)

type friendshipLayer struct {
	db *gorm.DB
}

type FriendshipRepo interface {
	Create(f *models.Friendship) error
	FindByUser(userID int) ([]models.Friendship, error)
	Find(userID, friendID int) (*models.Friendship, error)
	Update(f *models.Friendship) error
}

func newFriendshipLayer(db *gorm.DB) *friendshipLayer {
	return &friendshipLayer{
		db: db,
	}
}

func (l *friendshipLayer) Create(f *models.Friendship) error {
	return l.db.Create(f).Error
}

func (l *friendshipLayer) FindByUser(userID int) ([]models.Friendship, error) {
	var friends []models.Friendship
	err := l.db.Where("user_id = ?", userID).Find(&friends).Error
	return friends, err
}

func (l *friendshipLayer) Find(userID, friendID int) (*models.Friendship, error) {
	var f models.Friendship
	err := l.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&f).Error
	return &f, err
}

func (l *friendshipLayer) Update(f *models.Friendship) error {
	return l.db.Save(f).Error
}
