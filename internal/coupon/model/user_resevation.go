package model

import (
	"encoding/json"
)

type UserReservation struct {
	UserID         uint32 `json:"user_id" db:"user_id"`
	CouponActiveID uint64 `json:"coupon_active_id" db:"coupon_active_id"`
}

func (u UserReservation) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
