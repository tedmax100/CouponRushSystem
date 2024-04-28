package repository

import (
	"database/sql"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/tedmax100/CouponRushSystem/internal/types"
	"github.com/tedmax100/CouponRushSystem/internal/user/model"
)

type UserRepository struct {
	*sqlx.DB
	redisClient       *redis.Client
	couponActiveCache sync.Map //map[uint64]CacheItem
}

func NewUserRepository(db *sqlx.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		db,
		redis,
		sync.Map{},
	}
}

func (r *UserRepository) GetUser(id uint32) (model.User, error) {
	if item, exists := r.getFromCache(id); exists {
		return item, nil
	}

	user := model.User{}
	if err := r.DB.Get(&user, "SELECT id, name FROM user WHERE id = ? LIMIT 1", id); err != nil {
		if err == sql.ErrNoRows {
			return user, types.ErrorUserNotFound
		}
		return user, err
	}

	r.addToCache(id, user)

	return user, nil
}

type CacheItem struct {
	Value model.User
	Timer *time.Timer
}

func (r *UserRepository) addToCache(id uint32, item model.User) {
	now := time.Now()
	duration := now.Add(time.Hour * 1).Sub(now)

	timer := time.AfterFunc(duration, func() {
		r.removeFromCache(id)
	})

	r.couponActiveCache.Store(id, CacheItem{
		Value: item,
		Timer: timer,
	})

}

func (r *UserRepository) getFromCache(id uint32) (model.User, bool) {
	value, exist := r.couponActiveCache.Load(id)
	if exist {
		return value.(CacheItem).Value, true
	}
	return model.User{}, false
}

func (r *UserRepository) removeFromCache(id uint32) {
	r.couponActiveCache.Delete(id)
}
