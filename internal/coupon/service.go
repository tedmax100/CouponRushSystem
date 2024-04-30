package coupon

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/repository"
	//"github.com/tedmax100/CouponRushSystem/internal/message_queue"
)

type CouponActiveService struct {
	repo *repository.CouponActiveRepository
	//message_queue *message_queue.Broker
	resevedChan   chan<- model.UserReservedEvent
	purchasedChan chan<- model.PurchaseCouponEvent
}

func NewCouponActiveService(repo *repository.CouponActiveRepository, resevedChan chan model.UserReservedEvent, purchasedChan chan model.PurchaseCouponEvent) *CouponActiveService {
	return &CouponActiveService{
		repo:          repo,
		resevedChan:   resevedChan,
		purchasedChan: purchasedChan,
	}
}

func (s *CouponActiveService) GetActive(ctx context.Context, activeId uint64) (model.CouponActive, error) {
	return s.repo.GetActive(ctx, activeId)
}

func (s *CouponActiveService) ReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) error {
	resevedSeq, err := s.repo.ReserveCoupon(ctx, couponActive, userId)
	if err != nil {
		return err
	}

	err = s.ProcessAndGenerateCoupon(ctx, resevedSeq, couponActive)
	if err != nil {
		return err
	}

	event := model.UserReservedEvent{
		UserID:         userId,
		CouponActiveID: couponActive.ID,
	}

	select {
	case s.resevedChan <- event:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (s *CouponActiveService) PurchaseCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (model.Coupon, error) {
	_, err := s.repo.CheckReserveCoupon(ctx, couponActive, userId)
	if err != nil {
		return model.Coupon{}, err
	}

	purchasedCoupon, err := s.repo.PurchaseCoupon(ctx, couponActive, userId)
	if err != nil {
		return model.Coupon{}, err
	}

	event := model.PurchaseCouponEvent{
		UserID: userId,
		Coupon: purchasedCoupon,
	}

	select {
	case s.purchasedChan <- event:
	case <-ctx.Done():
		return model.Coupon{}, ctx.Err()
	}

	return purchasedCoupon, nil
}

func (s *CouponActiveService) ProcessAndGenerateCoupon(ctx context.Context, resevedSeq uint64, couponActive model.CouponActive) error {
	if resevedSeq%5 == 0 {
		couponCode, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		coupon := model.Coupon{
			Code:           couponCode.String(),
			CouponActiveID: couponActive.ID,
			CreatedAt:      time.Now().UTC(),
		}

		err = s.repo.AddCoupon(ctx, coupon)
		if err != nil {
			return err
		}
	}
	return nil
}
