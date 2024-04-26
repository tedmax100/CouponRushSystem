package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/tedmax100/CouponRushSystem/internal/config"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedisParams struct {
	fx.In

	Ctx    context.Context
	Config *config.Config
}

func NewRedis(p RedisParams) *redis.Client {
	log.Debug(context.Background(), "", zap.Any("config", p.Config))
	option, err := redis.ParseURL(p.Config.Redis)
	if err != nil {
		log.Fatal(context.Background(), err, zap.String("address", p.Config.Redis))
	}
	option.PoolSize = 5
	client := redis.NewClient(&redis.Options{
		Addr: option.Addr,
		//Username:  option.Username,
		//Password:  option.Password,
		PoolSize: option.PoolSize,
		Protocol: option.Protocol,
		//TLSConfig: option.TLSConfig,
	})
	if err := client.Ping(p.Ctx).Err(); err != nil {
		log.Fatal(context.Background(), err)
	}

	return client
}

func CloseRedis(client *redis.Client, lc fx.Lifecycle) {
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				err := client.Close()
				if err != nil {
					log.Error(ctx, err)
					return err
				}
				return nil
			},
		},
	)
}
