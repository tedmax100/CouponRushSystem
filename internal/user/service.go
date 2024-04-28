package user

import (
	"github.com/tedmax100/CouponRushSystem/internal/user/model"
	"github.com/tedmax100/CouponRushSystem/internal/user/repository"
)

type UserSertive struct {
	repo repository.UserRepository
}

func NewUserSertive(repo repository.UserRepository) *UserSertive {
	return &UserSertive{
		repo: repo,
	}
}

func (s *UserSertive) GetUser(id int32) (model.User, error) {
	return s.repo.GetUser(id)
}
