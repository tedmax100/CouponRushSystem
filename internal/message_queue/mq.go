package message_queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tedmax100/CouponRushSystem/internal/config"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	FANOUT                    = "fanout"
	DIRECT                    = "direct"
	TOPIC                     = "topic"
	NONE                      = ""
	TYPE_PROTOBUF             = "application/x-protobuf"
	TYPE_JSON                 = "application/json"
	ReconnectTimes            = 3
	ReconnectInterval         = 100
	ReservationActiveExchange = "ReservationActive"
	PurchaseCouponQueue       = "PurchaseCoupon"
)

type Broker struct {
	sync.Mutex
	ctx                   context.Context
	url                   string
	ch                    *amqp.Channel
	tokenEventCh          *amqp.Channel
	txQ                   amqp.Queue
	conn                  *amqp.Connection
	connected             bool
	close                 chan bool
	waitConnection        chan struct{}
	reservationActiveChan chan amqp.Delivery
	purchaseChan          chan amqp.Delivery
}

type BrokerParams struct {
	fx.In

	Ctx    context.Context
	Config config.Config
}

func NewBroker(p BrokerParams) *Broker {
	return &Broker{
		ctx:                   p.Ctx,
		url:                   p.Config.MQ,
		close:                 make(chan bool),
		waitConnection:        make(chan struct{}),
		reservationActiveChan: make(chan amqp.Delivery),
		purchaseChan:          make(chan amqp.Delivery),
	}
}

func (b *Broker) reconnect() {
	for {
		if b.connected {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if err := b.tryConnect(); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		b.Lock()
		b.connected = true
		close(b.waitConnection)
		b.channelSubsctibe()
		b.Unlock()

		notifyClose := make(chan *amqp.Error)
		b.conn.NotifyClose(notifyClose)

		chanNotifyClose := make(chan *amqp.Error)
		channelNotifyReturn := make(chan amqp.Return)

		if b.tokenEventCh != nil {
			channel := b.tokenEventCh
			channel.NotifyClose(chanNotifyClose)
			channel.NotifyReturn(channelNotifyReturn)
		}

		select {
		case result, ok := <-channelNotifyReturn:
			if ok {
				return
			}
			log.Error(b.ctx, fmt.Errorf("channelNotifyReturn"), zap.String("reason", result.ReplyText), zap.String("description", result.Exchange))
		case err := <-chanNotifyClose:
			log.Error(b.ctx, err, zap.String("msg", "chanNotifyClose"))
			b.setToDisconnect()
		case err := <-notifyClose:
			log.Error(b.ctx, err, zap.String("msg", "notifyClose"))
			select {
			case errs := <-chanNotifyClose:
				log.Error(b.ctx, errs)
			case <-time.After(time.Second):
			}
			b.setToDisconnect()
		case <-b.close:
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (b *Broker) connect() error {
	if err := b.tryConnect(); err != nil {
		log.Error(b.ctx, err)
		return err
	}
	b.Lock()
	b.connected = true
	b.Unlock()
	if err := b.Init(); err != nil {
		log.Error(b.ctx, err)
		return err
	}
	b.channelSubsctibe()
	go b.reconnect()

	return nil
}

func (b *Broker) tryConnect() error {
	amqpConfig := amqp.Config{
		Heartbeat: 10 * time.Second,
	}
	conn, err := amqp.DialConfig(b.url, amqpConfig)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to connect to RabbitMQ"))
		return (err)
	}
	b.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to create channel"))
		return err
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to set QOS"))
		return err
	}
	b.ch = ch

	exCh, err := conn.Channel()
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to create channel"))
		return err
	}
	err = exCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to set QOS"))
		return err
	}
	b.tokenEventCh = exCh
	return nil
}

func (b *Broker) Connect() error {
	b.Lock()
	if b.connected {
		b.Unlock()
		return nil
	}

	select {
	case <-b.close:
		b.close = make(chan bool)
	default:
	}
	b.Unlock()
	return b.connect()
}

func (b *Broker) Close() error {
	err := b.ch.Close()
	if err != nil {
		return err
	}
	// b.txCh.Close()
	err = b.tokenEventCh.Close()
	if err != nil {
		return err
	}
	err = b.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) Init() error {
	if err := b.DeclareExchanges(); err != nil {
		log.Error(b.ctx, err)
		return err
	}

	if err := b.DeclareQueues(); err != nil {
		log.Error(b.ctx, err)
		return err
	}
	return nil
}

func (b *Broker) DeclareExchanges() error {
	if err := b.ch.ExchangeDeclare(ReservationActiveExchange, FANOUT, true, false, false, false, nil); err != nil {
		log.Error(b.ctx, err)
		return err
	}
	return nil
}

func (b *Broker) DeclareQueues() error {
	_, err := b.ch.QueueDeclare(PurchaseCouponQueue, true, false, false, false, nil)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to declare a queue"))
		return err
	}
	err = b.ch.QueueBind(xxx, CANCEL_ORDER_PHASE, ROLLBACK_DISPATCH_EXCHANGE, false, nil)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to declare a queue"))
		return err
	}

	return nil
}

func (b *Broker) setToDisconnect() {
	b.Lock()
	b.connected = false
	b.waitConnection = make(chan struct{})
	b.Unlock()
}

func RunBroker(broker *Broker, lc fx.Lifecycle) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				err := broker.Connect()
				if err != nil {
					log.Fatal(broker.ctx, err)
					return err
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				err := broker.Close()
				if err != nil {
					log.Error(broker.ctx, err)
					return err
				}
				return nil
			},
		},
	)
}
