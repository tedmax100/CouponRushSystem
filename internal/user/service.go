package user

import (
	"github.com/tedmax100/CouponRushSystem/internal/user/model"
)

type UserSertive struct {
	repo UserRepository
}

func NewUserSertive(repo UserRepository) *UserSertive {
	return &UserSertive{
		repo: repo,
	}
}

func (s *UserSertive) GetUser(id uint64) (model.User, error) {
	return s.repo.GetUser(id)
}
