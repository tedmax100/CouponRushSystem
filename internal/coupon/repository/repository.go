package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
	"github.com/tedmax100/CouponRushSystem/internal/types"
)

type CouponActiveRepository struct {
	*sqlx.DB
	redisClient       *redis.Client
	couponActiveCache sync.Map
}

func NewActiveRepository(db *sqlx.DB, redis *redis.Client) *CouponActiveRepository {
	return &CouponActiveRepository{
		db,
		redis,
		sync.Map{},
	}
}

func (r *CouponActiveRepository) GetActive(ctx context.Context, activeId uint64) (model.CouponActive, error) {
	if item, exists := r.getFromCache(activeId); exists {
		return item, nil
	}

	couponActive := model.CouponActive{}
	if err := r.DB.GetContext(ctx, &couponActive, "SELECT id, date,begin_time, end_time, state begin FROM active WHERE id = ? LIMIT 1", activeId); err != nil {
		if err == sql.ErrNoRows {
			return couponActive, types.ErrorCouponActiveNotFound
		}

		return couponActive, err
	}
	return couponActive, nil
}

func (r *CouponActiveRepository) ReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (uint64, error) {
	dateKey := couponActive.ActiveDate.Format("coupan_active_20060102")
	reserveAtiveAmounKey := couponActive.ActiveDate.Format("reserve_active_amount_20060102")

	script := `
		local dateKey = KEYS[1]
		local reserveAtiveAmounKey = KEYS[2]
		local userId = ARGV[1]

		redis.call('SETBIT', dateKey, userId, 1)
		local reservedSeq = redis.call('INCR', reserveAtiveAmounKey)

		return reservedSeq
	`

	result, err := r.redisClient.Eval(ctx, script, []string{dateKey, reserveAtiveAmounKey}, userId).Result()
	if err != nil {
		return 0, err
	}

	resevedSeq, ok := result.(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type: %T, value: %v", result, result)
	}

	return resevedSeq, nil
}

func (r *CouponActiveRepository) AddCoupon(ctx context.Context, coupon model.Coupon) error {
	couponJSON, err := coupon.Marshal()
	if err != nil {
		return err
	}

	couponKey := "coupons_" + time.Now().Format("20060102")

	err = r.redisClient.LPush(ctx, couponKey, couponJSON).Err()
	if err != nil {
		return err
	}

	/*
	   // Create a coupon event
	*/
	return nil
}

func (r *CouponActiveRepository) CheckReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (bool, error) {
	dateKey := couponActive.ActiveDate.Format("coupan_active_20060102")

	bitValue, err := r.redisClient.GetBit(ctx, dateKey, int64(userId)).Result()
	if err != nil {
		return false, err
	}
	if bitValue != 1 {
		return false, types.ErrorUserNotReserveCouponActive
	}
	return true, nil
}

func (r *CouponActiveRepository) PurchaseCoupon(ctx context.Context, couponActive model.CouponActive, userId uint64) (model.Coupon, error) {
	couponPurchaseKey := couponActive.ActiveDate.Format("coupan_purchased_20060102")
	couponKey := "coupons_" + time.Now().Format("20060102")

	script := `
		if redis.call('GETBIT', KEYS[1], ARGV[1]) == 1 then
			return 'ALREADY_PURCHASED' -- User has already purchased a coupon, return specific message
		end

		-- Try to pop a coupon from the coupon list
		local coupon = redis.call('LPOP', KEYS[2])
		if coupon then
			-- Mark the user as having purchased a coupon
			redis.call('SETBIT', KEYS[1], ARGV[1], 1)
		
			return coupon
		else
			return 'NO_COUPONS'
		end
	`

	result, err := r.redisClient.Eval(ctx, script, []string{couponPurchaseKey, couponKey}, userId).Result()
	if err != nil {
		return model.Coupon{}, err
	}

	switch result.(string) {
	case "ALREADY_PURCHASED":
		return model.Coupon{}, types.ErrorUserAlreadyPurchasedCoupon
	case "NO_COUPONS":
		return model.Coupon{}, types.ErrorNoCouponToPurchase
	default:
		couponJSON, ok := result.(string)
		if !ok {
			return model.Coupon{}, fmt.Errorf("unexpected result type: %T, value: %v", result, result)
		}

		var coupon model.Coupon
		err = coupon.Unmarshal([]byte(couponJSON))
		if err != nil {
			return model.Coupon{}, fmt.Errorf("failed to unmarshal coupon: %v", err)
		}

		return coupon, nil
	}
}

type CacheItem struct {
	Value model.CouponActive
	Timer *time.Timer
}

func (r *CouponActiveRepository) addToCache(id uint64, item model.CouponActive) {
	duration := item.ActiveEndTime.Sub(time.Now())
	timer := time.AfterFunc(duration, func() {
		r.removeFromCache(id)
	})
	r.couponActiveCache.Store(id, CacheItem{
		Value: item,
		Timer: timer,
	})
}

func (r *CouponActiveRepository) getFromCache(id uint64) (model.CouponActive, bool) {
	value, exist := r.couponActiveCache.Load(id)
	if exist {
		return value.(CacheItem).Value, true
	}
	return model.CouponActive{}, false
}

func (r *CouponActiveRepository) removeFromCache(id uint64) {
	r.couponActiveCache.Delete(id)
}
