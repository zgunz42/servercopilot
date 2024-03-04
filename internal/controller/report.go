package controller

import (
	"bufio"
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/reactivex/rxgo/v2"
	"github.com/valyala/fasthttp"
	"github.com/zgunz42/servercopilot/internal/stream"
	"go.uber.org/fx"
)

type ReportController struct {
	fx.In
	TempSensor *stream.TempSensor
	HumSensor  *stream.HumSensor
	Route      fiber.Router `name:"api-v1"`
}

func (c *ReportController) GetSensor(ctx *fiber.Ctx) error {
	humObs := c.HumSensor.GetHumObs().Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		return map[string]any{
			"hum": i,
		}, nil
	})
	tempObs := c.TempSensor.GetTempObs().Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		return map[string]any{
			"temp": i,
		}, nil
	})
	sensors := []rxgo.Observable{humObs, tempObs}
	// zip the two observables
	obs := rxgo.Merge(sensors, rxgo.WithBufferedChannel(2))

	_ctx := ctx.Context()

	_ctx.SetContentType("text/event-stream")
	_ctx.Response.Header.Set("Cache-Control", "no-cache")
	_ctx.Response.Header.Set("Connection", "keep-alive")
	_ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	_ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	_ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	_ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

	_, cancel := obs.Connect(_ctx)

	go func() {
		log.Debug("cancelling report controller")

		cancel()
	}()

	_ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case <-_ctx.Done():
				log.Debug("finished report controller")

				w.Flush()
				return
			case val := <-obs.Observe():
				data := val.V.(map[string]any)
				encData, err := json.Marshal(data)

				if err != nil {
					panic(err)
				}

				fmt.Fprintf(w, "data: %s\n\n", encData)
				w.Flush()
			}
		}
	}))

	return nil
}

func (c *ReportController) GetHum(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]any{
		"hum": c.HumSensor.GetHum(ctx.Context()),
	})
}

func (c *ReportController) GetTemp(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]any{
		"temp": c.TempSensor.GetTemp(ctx.Context()),
	})
}

func RegisterReportController(c ReportController) {
	c.Route.Get("/temp", c.GetTemp)
	c.Route.Get("/hum", c.GetHum)
	c.Route.Get("/watch", c.GetSensor)
}
