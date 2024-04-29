package coupon

import (
	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/repository"
	"github.com/tedmax100/CouponRushSystem/internal/message_queue"
)

type CouponActiveReceiver struct {
	repo          *repository.CouponActiveRepository
	message_queue *message_queue.Broker
	resevedChan   <-chan model.UserReservedEvent
	purchasedChan <-chan model.PurchaseCouponEvent
}

func NewCouponActiveReceiver(repo *repository.CouponActiveRepository, mq *message_queue.Broker, resevedChan chan model.UserReservedEvent, purchasedChan chan model.PurchaseCouponEvent) *CouponActiveReceiver {
	return &CouponActiveReceiver{
		repo:          repo,
		message_queue: mq,
		resevedChan:   resevedChan,
		purchasedChan: purchasedChan,
	}
}
