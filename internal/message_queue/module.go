package message_queue

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewBroker),
	fx.Invoke(RunBroker),
)
