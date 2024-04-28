package model

import (
	"encoding/json"
	"time"
)

type Coupon struct {
	//ID             uint64    `json:"id" db:"id"`
	Code           string    `json:"code" db:"code"`
	CouponActiveID uint64    `json:"activeId" db:"activeId"`
	CreatedAt      time.Time `json:"createdAt" db:"createdAt"`
}

func (u Coupon) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Coupon) Unmarshal(data []byte) error {
	return json.Unmarshal(data, u)
}

type PurchaseCouponEvent struct {
	UserID uint32 `json:"userId"`
	Coupon
}
