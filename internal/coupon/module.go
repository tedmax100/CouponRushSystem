package coupon

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewCouponActiveService),
	fx.Provide(NewCouponEventReceiverService),
	fx.Invoke(RunReceiver),
)
