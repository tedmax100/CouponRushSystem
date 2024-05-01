package coupon

import (
	"context"

	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
)

type CouponActiveRepositoryInterface interface {
	GetActive(ctx context.Context, activeId uint64) (model.CouponActive, error)
	ReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (uint64, error)
	CheckReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (bool, error)
	PurchaseCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (model.Coupon, error)
	AddCoupon(ctx context.Context, coupon model.Coupon) error
}
