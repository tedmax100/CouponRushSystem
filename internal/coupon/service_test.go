package coupon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
)

type mockCouponActiveRepository struct {
	CouponActiveRepositoryInterface
	checkReserveCouponErr error
	purchaseCouponErr     error
}

func (m *mockCouponActiveRepository) CheckReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (bool, error) {
	return true, m.checkReserveCouponErr
}

func (m *mockCouponActiveRepository) GetActive(ctx context.Context, activeId uint64) (model.CouponActive, error) {
	return model.CouponActive{ID: activeId}, nil
}

func (m *mockCouponActiveRepository) ReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (uint64, error) {
	if couponActive.State != model.OPENING {
		return 0, errors.New("coupon active is not open")
	}

	now := time.Now()
	if now.Before(couponActive.ActiveBeginTime) || now.After(couponActive.ActiveEndTime) {
		return 0, errors.New("current time is not within the coupon active's start and end times")
	}

	return 1, nil
}

func (m *mockCouponActiveRepository) PurchaseCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (model.Coupon, error) {
	if m.purchaseCouponErr != nil {
		return model.Coupon{}, m.purchaseCouponErr
	}
	couponCode, _ := uuid.NewRandom()

	return model.Coupon{Code: couponCode.String(), CouponActiveID: couponActive.ID}, nil
}

func TestGetActive(t *testing.T) {
	// Arrange
	repo := &mockCouponActiveRepository{}
	service := NewCouponActiveService(repo, make(chan model.UserReservedEvent), make(chan model.PurchaseCouponEvent))

	//Act
	active, err := service.GetActive(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), active.ID)
}

func TestReserveCoupon(t *testing.T) {
	ctx := context.Background()

	userId := uint64(1)

	t.Run("returns error when ReserveCoupon fails", func(t *testing.T) {
		// Arrange
		couponActive := model.CouponActive{ID: 1, State: model.NOT_OPEN, ActiveDate: time.Now().Add(time.Hour), ActiveBeginTime: time.Now(), ActiveEndTime: time.Now().Add(2 * time.Hour)}
		repo := &mockCouponActiveRepository{checkReserveCouponErr: errors.New("coupon active is not open")}
		reservedChan := make(chan model.UserReservedEvent)
		go func() {
			for range reservedChan {

			}
		}()

		// Act
		service := NewCouponActiveService(repo, reservedChan, make(chan model.PurchaseCouponEvent))

		// Assert
		err := service.ReserveCoupon(ctx, couponActive, userId)
		assert.EqualError(t, err, "coupon active is not open")
	})

	t.Run("returns no error when successful", func(t *testing.T) {
		// Arrange
		couponActive := model.CouponActive{ID: 1, State: model.OPENING, ActiveDate: time.Now(), ActiveBeginTime: time.Now(), ActiveEndTime: time.Now().Add(time.Hour)}
		repo := &mockCouponActiveRepository{}
		reservedChan := make(chan model.UserReservedEvent)
		go func() {
			for range reservedChan {

			}
		}()

		// Act
		service := NewCouponActiveService(repo, reservedChan, make(chan model.PurchaseCouponEvent))

		// Assert
		err := service.ReserveCoupon(ctx, couponActive, userId)
		assert.NoError(t, err)
	})
}

func TestPurchaseCoupon(t *testing.T) {
	ctx := context.Background()
	couponActive := model.CouponActive{ID: 1}
	userId := uint64(1)

	t.Run("returns error when CheckReserveCoupon fails", func(t *testing.T) {
		// Arrange
		repo := &mockCouponActiveRepository{checkReserveCouponErr: errors.New("check reserve coupon error")}

		service := NewCouponActiveService(repo, make(chan model.UserReservedEvent), make(chan model.PurchaseCouponEvent))

		// Act
		_, err := service.PurchaseCoupon(ctx, couponActive, userId)

		// Assert
		assert.EqualError(t, err, "check reserve coupon error")
	})

	t.Run("returns error when PurchaseCoupon fails", func(t *testing.T) {
		// Arrange
		repo := &mockCouponActiveRepository{purchaseCouponErr: errors.New("purchase coupon error")}
		service := NewCouponActiveService(repo, make(chan model.UserReservedEvent), make(chan model.PurchaseCouponEvent))

		// Act
		_, err := service.PurchaseCoupon(ctx, couponActive, userId)

		// Assert
		assert.EqualError(t, err, "purchase coupon error")
	})

	t.Run("returns purchased coupon when successful", func(t *testing.T) {
		// Arrange
		repo := &mockCouponActiveRepository{}
		purchasedChan := make(chan model.PurchaseCouponEvent)
		go func() {
			for range purchasedChan {

			}
		}()

		service := NewCouponActiveService(repo, make(chan model.UserReservedEvent), purchasedChan)

		// Act
		coupon, err := service.PurchaseCoupon(ctx, couponActive, userId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, couponActive.ID, coupon.CouponActiveID)
	})
}
