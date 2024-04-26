package database

import (
	"context"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tedmax100/CouponRushSystem/internal/config"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	POSTGRES string = "postgres"
)

type DbParams struct {
	fx.In

	Ctx    context.Context
	Config *config.Config
}

func NewDB(p DbParams) *sqlx.DB {
	dsn := p.Config.DB
	dsn += "?application_name=coupon_rush_server&sslmode=disable"
	oteldb, err := otelsql.Open(POSTGRES, dsn)
	if err != nil {
		log.Fatal(context.Background(), err)
	}
	err = otelsql.RegisterDBStatsMetrics(oteldb, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		log.Fatal(context.Background(), err)
	}

	db := sqlx.NewDb(oteldb, POSTGRES)

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatal(context.Background(), err, zap.Any("msg", "Database init failed"))
	}

	return db
}

func CloseDb(client *sqlx.DB, lc fx.Lifecycle) {
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
