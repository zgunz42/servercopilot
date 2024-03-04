package stream

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("stm", fx.Provide(NewPubSub), fx.Provide(CreateTempSens), fx.Provide(CreateHumSens), fx.Provide())
}
