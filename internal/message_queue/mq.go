package message_queue

import (
	"context"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tedmax100/CouponRushSystem/internal/config"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	FANOUT                       = "fanout"
	DIRECT                       = "direct"
	TOPIC                        = "topic"
	NONE                         = ""
	TYPE_PROTOBUF                = "application/x-protobuf"
	TYPE_JSON                    = "application/json"
	ReconnectTimes               = 3
	ReconnectInterval            = 100
	CouponEventExchange          = "CouponEvent"
	UserReserveCouponActiveQueue = "UserReserveCouponQueue"
	PurchaseCouponQueue          = "PurchaseCouponQueue"
	UserRserveCouponActivekey    = "UserReserveCouponActive"
	PurchaseCouponKey            = "PurchaseCoupon"
)

type BrokerParams struct {
	fx.In

	Ctx    context.Context
	Config *config.Config
}

type Broker struct {
	sync.Mutex
	ctx            context.Context
	url            string
	conn           *amqp.Connection
	connected      bool
	close          chan bool
	waitConnection chan struct{}
	ch             *amqp.Channel
	//newCouponActiveChan *amqp.Channel
	//userReseveCouponActiveChan *amqp.Channel
	//purchaseChan               *amqp.Channel
}

func NewBroker(p BrokerParams) *Broker {
	return &Broker{
		ctx:            p.Ctx,
		url:            p.Config.MQ,
		close:          make(chan bool),
		waitConnection: make(chan struct{}),
		//reservationActiveChan: make(chan amqp.Delivery),
		//purchaseChan:          make(chan amqp.Delivery),
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
		b.Unlock()

		notifyClose := make(chan *amqp.Error)
		b.conn.NotifyClose(notifyClose)

		select {

		case err := <-notifyClose:
			log.Error(b.ctx, err, zap.String("msg", "notifyClose"))
			select {
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
	if err := b.initial(); err != nil {
		log.Error(b.ctx, err)
		return err
	}
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
		5,
		0,
		false,
	)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to set QOS"))
		return err
	}
	b.ch = ch

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

func (b *Broker) closeConn() error {
	err := b.ch.Close()
	if err != nil {
		return err
	}
	err = b.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) initial() error {
	if err := b.declareExchanges(); err != nil {
		log.Error(b.ctx, err)
		return err
	}

	if err := b.declareQueues(); err != nil {
		log.Error(b.ctx, err)
		return err
	}
	return nil
}

func (b *Broker) declareExchanges() error {
	return b.ch.ExchangeDeclare(
		CouponEventExchange,
		DIRECT,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (b *Broker) declareQueues() error {
	_, err := b.ch.QueueDeclare(
		UserReserveCouponActiveQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = b.ch.QueueDeclare(
		PurchaseCouponQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err = b.ch.QueueBind(UserReserveCouponActiveQueue, UserRserveCouponActivekey, CouponEventExchange, false, nil); err != nil {
		return err
	}

	if err = b.ch.QueueBind(PurchaseCouponQueue, PurchaseCouponKey, CouponEventExchange, false, nil); err != nil {
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

func (b *Broker) SendReservedCouponActiveEvent(ctx context.Context, msg []byte) error {
	err := b.ch.PublishWithContext(
		ctx,
		CouponEventExchange,
		UserRserveCouponActivekey,
		false,
		false,
		amqp.Publishing{
			ContentType: TYPE_JSON,
			Body:        msg,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) SendPurchaseCouponEvent(ctx context.Context, msg []byte) error {
	err := b.ch.PublishWithContext(
		ctx,
		CouponEventExchange,
		PurchaseCouponKey,
		false,
		false,
		amqp.Publishing{
			ContentType: TYPE_JSON,
			Body:        msg,
		},
	)
	if err != nil {
		return err
	}
	return nil
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
				err := broker.closeConn()
				if err != nil {
					log.Error(broker.ctx, err)
					return err
				}
				return nil
			},
		},
	)
}
