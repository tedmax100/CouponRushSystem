package repository

import (
	"context"
	"database/sql"
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

func (r *CouponActiveRepository) ReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint32) error {
	dateKey := couponActive.ActiveDate.Format("coupan_active_20060102")

	err := r.redisClient.SetBit(ctx, dateKey, int64(userId), 1).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *CouponActiveRepository) CheckReserveCoupon(ctx context.Context, couponActive model.CouponActive, userId uint32) (bool, error) {
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
