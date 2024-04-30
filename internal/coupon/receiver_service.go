package coupon

import (
	"context"
	"sync"

	"github.com/tedmax100/CouponRushSystem/internal/log"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
	"github.com/tedmax100/CouponRushSystem/internal/message_queue"
)

type CouponEventReceiverService struct {
	//repo          *repository.CouponActiveRepository
	message_queue *message_queue.Broker
	resevedChan   <-chan model.UserReservedEvent
	purchasedChan <-chan model.PurchaseCouponEvent
}

func NewCouponEventReceiverService(
	mq *message_queue.Broker,
	resevedChan chan model.UserReservedEvent,
	purchasedChan chan model.PurchaseCouponEvent) *CouponEventReceiverService {
	return &CouponEventReceiverService{
		//repo:          repo,
		message_queue: mq,
		resevedChan:   resevedChan,
		purchasedChan: purchasedChan,
	}
}

func RunReceiver(ctx context.Context, r *CouponEventReceiverService, lc fx.Lifecycle) {
	var wg sync.WaitGroup

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				wg.Add(2)

				go func() {
					defer wg.Done()

					for {
						select {
						case resevedEvent, ok := <-r.resevedChan:
							if !ok {
								return
							}

							msgBytes, err := resevedEvent.Marshal()
							if err != nil {
								log.Error(context.Background(), err, zap.Any("event", resevedEvent))
								continue
							}
							r.message_queue.SendReservedCouponActiveEvent(context.Background(), msgBytes)
						case <-ctx.Done():
							return
						}
					}
				}()

				go func() {
					defer wg.Done()

					for {
						select {
						case purchasedEvent, ok := <-r.purchasedChan:
							if !ok {
								return
							}

							msgBytes, err := purchasedEvent.Marshal()
							if err != nil {
								log.Error(context.Background(), err, zap.Any("event", purchasedEvent))
								continue
							}
							r.message_queue.SendPurchaseCouponEvent(context.Background(), msgBytes)
						case <-ctx.Done():
							return
						}
					}
				}()

				return nil
			},
			OnStop: func(context.Context) error {
				// Wait for all goroutines to finish
				wg.Wait()

				// Consume all remaining events in the channels
				for resevedEvent := range r.resevedChan {
					msgBytes, err := resevedEvent.Marshal()
					if err != nil {
						log.Error(context.Background(), err, zap.Any("event", resevedEvent))
						continue
					}
					r.message_queue.SendReservedCouponActiveEvent(context.Background(), msgBytes)
				}

				for purchasedEvent := range r.purchasedChan {
					msgBytes, err := purchasedEvent.Marshal()
					if err != nil {
						log.Error(context.Background(), err, zap.Any("event", purchasedEvent))
						continue
					}
					r.message_queue.SendPurchaseCouponEvent(context.Background(), msgBytes)
				}

				return nil
			},
		},
	)
}
