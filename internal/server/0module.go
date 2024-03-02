package server

import (
	"go.uber.org/fx"

	"github.com/zgunz42/servercopilot/internal/server/http"
)

func Module() fx.Option {
	return fx.Module("server", http.Module())
}
