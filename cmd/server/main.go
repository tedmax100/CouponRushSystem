package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tedmax100/CouponRushSystem/pkg/log"
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
		fx.Logger(&log.FxLogger{}),
		fx.StopTimeout(time.Minute),
	)

	app.Run()
}

func init() {
	log.Initital(GitCommit, GitBranch, GoVersion, BuildDate)

}
