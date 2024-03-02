package cli

import (
	"context"

	"go.uber.org/fx"

	"github.com/zgunz42/servercopilot/internal/app"
	"github.com/zgunz42/servercopilot/internal/app/appcontext"
)

func Start(module fx.Option) {
	app.New(appcontext.Declare(appcontext.EnvCLI), module).Start(context.Background())
}
