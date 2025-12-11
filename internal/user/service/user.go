package service

import (
	"cashapp/core"
	"cashapp/internal/user/repository"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	repository repository.Repo
	config     *core.Config
}

func New(r repository.Repo, c *core.Config) *UserService {
	return &UserService{
		repository: r,
		config:     c,
	}
}

func (s *UserService) CreateUser(req core.CreateUserRequest) core.Response {
	user, err := s.repository.Users.FindByTag(req.Tag)

	if err == nil {
		return core.Error(errors.New("cash tag taken"), core.String("cash tag has already been taken"))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return core.Error(err, nil)
	}

	if err := s.repository.Users.Create(user); err != nil {
		return core.Error(err, nil)
	}

	wallet, err := s.repository.Wallets.Create(user.ID)
	if err != nil {
		return core.Error(err, nil)
	}

	user.Wallets = append(user.Wallets, *wallet)
	return core.Success(&map[string]interface{}{
		"user": user,
	}, core.String("user created successfully"))
}
