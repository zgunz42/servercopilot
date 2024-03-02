package server

import (
	"go.uber.org/fx"

	"github.com/zgunz42/servercopilot/internal/server/http"
	"github.com/zgunz42/servercopilot/internal/server/mqtt"
)

func Module() fx.Option {
	return fx.Module("server", http.Module(), mqtt.Module())
}
