package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"github.com/tedmax100/CouponRushSystem/internal/config"
	"github.com/tedmax100/CouponRushSystem/internal/database"
	"github.com/tedmax100/CouponRushSystem/internal/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	GitBranch = ""
	GoVersion = ""
	GitCommit = ""
	BuildDate = ""
)

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
		config.Module,
		database.Module,
		fx.Invoke(func(cancelFunc context.CancelFunc, lc fx.Lifecycle) {
			lc.Append(
				fx.Hook{
					OnStop: func(ctx context.Context) error {
						log.Info(ctx, "Process Exit")
						cancelFunc()
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

func runServer(configObj *config.Config, lc fx.Lifecycle) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.ContextWithFallback = true
	r.Use(otelgin.Middleware("coupun_rush_server"))
	//api.SetupRouter(r, GitCommit, baseRepo, mainRepo, settingRepo, jobRepo, contract, bondContract,

	server := &http.Server{
		Addr:    ":" + configObj.Port,
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
