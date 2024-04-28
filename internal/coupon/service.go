package coupon

import (
	"context"

	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/repository"
	"github.com/tedmax100/CouponRushSystem/internal/message_queue"
)

type CouponActiveService struct {
	repo          repository.CouponActiveRepository
	message_queue *message_queue.Broker
}

func NewCouponActiveService(repo repository.CouponActiveRepository) *CouponActiveService {
	return &CouponActiveService{
		repo: repo,
	}
}

func (s *CouponActiveService) GetActive(ctx context.Context, activeId uint64) (model.CouponActive, error) {
	return s.repo.GetActive(ctx, activeId)
}

func (s *CouponActiveService) ReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint32) error {
	if err := s.repo.ReserveCoupon(ctx, couponActive, userId); err != nil {
		return err
	}

	event := model.UserReservation{
		UserID:         userId,
		CouponActiveID: couponActive.ID,
	}

	binaryBody, err := event.Marshal()
	if err != nil {
		return err
	}
	if err := s.message_queue.SendPurchaseCouponEvent(ctx, binaryBody); err != nil {
		return err
	}

	return nil
}

func (s *CouponActiveService) PurchaseCoupon(ctx context.Context, couponActive model.CouponActive, userId uint32) error {
	_, err := s.repo.CheckReserveCoupon(ctx, couponActive, userId)
	if err != nil {
		return err
	}

	if err := s.repo.PurchaseCoupon(ctx, couponActive, userId); err != nil {
		return err
	}

	event := model.UserReservation{
		UserID:         userId,
		CouponActiveID: couponActive.ID,
	}

	binaryBody, err := event.Marshal()
	if err != nil {
		return err
	}
	if err := s.message_queue.SendPurchaseCouponEvent(ctx, binaryBody); err != nil {
		return err
	}

	return nil
}
