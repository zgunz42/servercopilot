package server

import (
	"context"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"github.com/zgunz42/servercopilot/internal/app"
	"github.com/zgunz42/servercopilot/internal/app/appconfig"
	"github.com/zgunz42/servercopilot/internal/app/appcontext"
	"github.com/zgunz42/servercopilot/internal/server/mqtt"
	"github.com/zgunz42/servercopilot/internal/stream"
)

func Run() {
	app.New(appcontext.Declare(appcontext.EnvServer), fx.Invoke(run)).Run()
}

func run(lc fx.Lifecycle, app *fiber.App, mqtt *mqtt.MqttClient, temp *stream.TempSensor, hum *stream.HumSensor, conf *appconfig.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := mqtt.Connect(3)

			if err != nil {
				return err
			}

			temp.Sub()
			hum.Sub()

			ln, err := net.Listen("tcp", conf.ServiceListenAddress)
			if err != nil {
				return err
			}

			go func() {
				if err := app.Listener(ln); err != nil {
					log.Error().Err(err).Msg("server terminated unexpectedly")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := mqtt.Disconnect()
			if err != nil {
				return err
			}

			err = hum.Close()
			if err != nil {
				return err
			}

			err = temp.Close()
			if err != nil {
				return err
			}

			log.Info().Msg("gracefully shutting down server")
			if err := app.Shutdown(); err != nil {
				log.Error().Err(err).Msg("error occurred while gracefully shutting down server")
				return err
			}
			log.Info().Msg("graceful server shut down completed")
			return nil
		},
	})
}
