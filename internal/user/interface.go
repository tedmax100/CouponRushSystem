package user

import (
	"github.com/tedmax100/CouponRushSystem/internal/user/model"
)

type UserRepository interface {
	GetUser(id uint64) (model.User, error)
}
