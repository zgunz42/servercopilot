package http

import (
	"go.uber.org/fx"

	"github.com/zgunz42/servercopilot/internal/server/http/route"
)

func Module() fx.Option {
	return fx.Module("http", fx.Provide(Create), route.Module())
}
