package controller

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

	_ctx := ctx.Context()

	_ctx.SetContentType("text/event-stream")
	_ctx.Response.Header.Set("Cache-Control", "no-cache")
	_ctx.Response.Header.Set("Connection", "keep-alive")
	_ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	_ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	_ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	_ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

	_ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case <-_ctx.Done():
				log.Debug("finished report controller")

				w.Flush()
				return
			case val := <-c.HumSensor.GetMsg():
				res, err := strconv.ParseFloat(string(val.Payload), 64)

				if err != nil {
					val.Ack()
					continue
				}

				encData, err := json.Marshal(map[string]any{
					"hum": res,
				})

				if err != nil {
					panic(err)
				}
				val.Ack()
				fmt.Fprintf(w, "data: %s\n\n", encData)
				w.Flush()
			case val := <-c.TempSensor.GetMsg():
				res, err := strconv.ParseFloat(string(val.Payload), 64)

				if err != nil {
					val.Ack()
					continue
				}

				encData, err := json.Marshal(map[string]any{
					"temp": res,
				})

				if err != nil {
					panic(err)
				}

				val.Ack()
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
