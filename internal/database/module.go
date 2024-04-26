package database

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewDB),
	fx.Provide(NewRedis),
	fx.Invoke(CloseRedis),
	fx.Invoke(CloseDb),
)
