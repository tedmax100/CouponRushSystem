package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"github.com/tedmax100/CouponRushSystem/internal/api"
	"github.com/tedmax100/CouponRushSystem/internal/config"
	"github.com/tedmax100/CouponRushSystem/internal/coupon"
	"github.com/tedmax100/CouponRushSystem/internal/coupon/model"
	couponRepo "github.com/tedmax100/CouponRushSystem/internal/coupon/repository"
	"github.com/tedmax100/CouponRushSystem/internal/database"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	"github.com/tedmax100/CouponRushSystem/internal/message_queue"
	"github.com/tedmax100/CouponRushSystem/internal/user"
	userRepo "github.com/tedmax100/CouponRushSystem/internal/user/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	GitBranch = ""
	GoVersion = ""
	GitCommit = ""
	BuildDate = ""
)

//	@title			Swagger Coupon Rush Server API
//  @version 1.0
//	@description	This is the Coupon Rush ServerOpenAPI.
//	@contact.name	Developer Support - Nathan
//	@contact.email	tedmax100@gmail.com
//	@license.name	Proprietary
//  @externalDocs.description Coupon Rush Server Architecture Documentation

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error(context.Background(), fmt.Errorf("panic"), zap.Any("recover", r))
		}
		log.Sync()
	}()

	ctx, cancel := context.WithCancel(context.TODO())

	app := fx.New(
		fx.Provide(func() context.Context {
			return ctx
		}),
		fx.Provide(func() context.CancelFunc {
			return cancel
		}),
		fx.Provide(func() chan model.UserReservedEvent {
			return make(chan model.UserReservedEvent)
		}),
		fx.Provide(func() chan model.PurchaseCouponEvent {
			return make(chan model.PurchaseCouponEvent)
		}),
		fx.Provide(func() *gin.Engine {
			r := gin.New()
			r.Use(gin.Recovery())
			r.ContextWithFallback = true
			r.Use(otelgin.Middleware("coupun_rush_server"))
			return r
		}),
		config.Module,
		database.Module,
		message_queue.Module,
		coupon.Module,
		couponRepo.Module,
		user.Module,
		userRepo.Module,
		fx.Invoke(func(cancelFunc context.CancelFunc, lc fx.Lifecycle, reservedChan chan model.UserReservedEvent, purchasedChan chan model.PurchaseCouponEvent) {
			lc.Append(
				fx.Hook{
					OnStop: func(ctx context.Context) error {
						log.Info(ctx, "Process Exit")
						cancelFunc()

						// Close the channels
						close(reservedChan)
						close(purchasedChan)

						return nil
					},
				},
			)
		}),
		fx.Invoke(runServer),
		fx.Logger(&log.FxLogger{}),
		fx.StopTimeout(time.Minute),
	)

	app.Run()
}

type ServerParams struct {
	fx.In
	CouponService *coupon.CouponActiveService
	UserService   *user.UserSertive
	HandlerParams api.HandlerParams
	ConfigObj     *config.Config
}

func runServer(p ServerParams, lc fx.Lifecycle) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.ContextWithFallback = true
	r.Use(otelgin.Middleware("coupun_rush_server"))

	api.SetupRouter(r, p.HandlerParams)

	server := &http.Server{
		Addr:    ":" + p.ConfigObj.Port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info(ctx, "service terminating")
			return server.Shutdown(ctx)
		},
	})
}

func init() {
	log.Initital(GitCommit, GitBranch, GoVersion, BuildDate)
}
