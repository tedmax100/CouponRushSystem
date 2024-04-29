package dto

import (
	"time"
)

type CommonErrorResponse struct {
	Code    int    `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ReserveCouponRequest struct {
	ActiveID uint64 `json:"active_id" binding:"required"`
	UserID   uint64 `json:"user_id" binding:"required"`
}
type ReserveCouponResponse struct {
	UserID       uint64    `json:"user_id"`
	ActiveID     uint64    `json:"active_id"`
	ReservedTime time.Time `json:"reserved_time"`
}

type PurchaseCouponRequest struct {
	ActiveID uint64 `json:"active_id" binding:"required"`
	UserID   uint64 `json:"user_id" binding:"required"`
}

type PurchaseCouponResponse struct {
	UserID     uint64 `json:"user_id"`
	ActiveID   uint64 `json:"active_id"`
	CouponCode string `json:"coupon_code"`
}
