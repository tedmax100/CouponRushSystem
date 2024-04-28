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
	ReservationActiveExchange    = "ReservationActiveExchange"
	UserReserveCouponActiveQueue = "UserReserveCouponQueue"
	PurchaseCouponQueue          = "PurchaseCouponQueue"
)

type Broker struct {
	sync.Mutex
	ctx                        context.Context
	url                        string
	conn                       *amqp.Connection
	connected                  bool
	close                      chan bool
	waitConnection             chan struct{}
	reservationActiveChan      *amqp.Channel
	userReseveCouponActiveChan *amqp.Channel
	purchaseChan               *amqp.Channel
}

type BrokerParams struct {
	fx.In

	Ctx    context.Context
	Config config.Config
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
	if err := b.Init(); err != nil {
		log.Error(b.ctx, err)
		return err
	}
	//b.channelSubsctibe()
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

	reservationActiveChan, err := conn.Channel()
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to create channel"))
		return err
	}
	err = reservationActiveChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to set QOS"))
		return err
	}
	b.reservationActiveChan = reservationActiveChan

	userReseveCouponActiveChan, err := conn.Channel()
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to create channel"))
		return err
	}
	err = userReseveCouponActiveChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to set QOS"))
		return err
	}
	b.userReseveCouponActiveChan = userReseveCouponActiveChan

	purchaseChan, err := conn.Channel()
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to create channel"))
		return err
	}
	err = purchaseChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error(b.ctx, err, zap.String("msg", "Failed to set QOS"))
		return err
	}
	b.purchaseChan = purchaseChan
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
	err := b.reservationActiveChan.Close()
	if err != nil {
		return err
	}
	err = b.purchaseChan.Close()
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
	err := b.reservationActiveChan.ExchangeDeclare(
		ReservationActiveExchange,
		FANOUT,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) DeclareQueues() error {
	_, err := b.reservationActiveChan.QueueDeclare(
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
	err = b.reservationActiveChan.QueueBind(
		PurchaseCouponQueue,
		"",
		ReservationActiveExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = b.purchaseChan.QueueDeclare(
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

	_, err = b.userReseveCouponActiveChan.QueueDeclare(
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
	return nil
}

func (b *Broker) setToDisconnect() {
	b.Lock()
	b.connected = false
	b.waitConnection = make(chan struct{})
	b.Unlock()
}

func (b *Broker) SendPurchaseCouponEvent(ctx context.Context, msg []byte) error {
	err := b.userReseveCouponActiveChan.PublishWithContext(
		ctx,
		"",
		UserReserveCouponActiveQueue,
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
