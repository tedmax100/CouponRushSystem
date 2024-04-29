package model

import (
	"encoding/json"
)

type UserReservedEvent struct {
	UserID         uint64 `json:"user_id" db:"user_id"`
	CouponActiveID uint64 `json:"coupon_active_id" db:"coupon_active_id"`
}

func (u UserReservedEvent) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
