package app

import (
	"go.uber.org/fx"

	"github.com/zgunz42/servercopilot/internal/app/appconfig"
	"github.com/zgunz42/servercopilot/internal/app/appcontext"
	"github.com/zgunz42/servercopilot/internal/controller"
	"github.com/zgunz42/servercopilot/internal/server"
	"github.com/zgunz42/servercopilot/internal/x/logger"
	"github.com/zgunz42/servercopilot/internal/x/logger/fxlogger"
)

func New(ctx appcontext.Ctx, additionalOpts ...fx.Option) *fx.App {
	conf, err := appconfig.Parse(ctx)
	if err != nil {
		panic(err)
	}

	// logger and configuration are the only two things that are not in the fx graph
	// because some other packages need them to be initialized before fx starts
	logger.Configure(conf)

	baseOpts := []fx.Option{
		fx.WithLogger(fxlogger.Logger),
		fx.Supply(conf),
		controller.Module(),
		server.Module(),
	}

	return fx.New(append(baseOpts, additionalOpts...)...)
}
