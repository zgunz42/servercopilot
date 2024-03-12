package controller

import (
	"bytes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/zgunz42/servercopilot/internal/model"
	"github.com/zgunz42/servercopilot/lib"
	"go.uber.org/fx"
)

type DeviceController struct {
	fx.In
	Route fiber.Router `name:"api-v1"`
}

func (c *DeviceController) GetVersion(ctx *fiber.Ctx) error {
	versionUrl := os.Getenv("GITHUB_RELEASE_API")
	result, err := lib.GetUrl[model.GithubRelease](versionUrl, model.Unmarshal)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(map[string]any{
		"version":    result.TagName,
		"url":        result.HTMLURL,
		"draft":      result.Draft,
		"prerelease": result.Prerelease,
	})

}

func (c *DeviceController) GetFirmware(ctx *fiber.Ctx) error {
	versionUrl := os.Getenv("GITHUB_RELEASE_API")
	result, err := lib.GetUrl[model.GithubRelease](versionUrl, model.Unmarshal)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	for _, asset := range result.Assets {
		if asset.Name == "firmware.bin" {
			downloadUrl := asset.BrowserDownloadURL
			result, err := lib.DownloadFile(downloadUrl, true)

			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error": err.Error(),
				})
			}

			reader := bytes.NewReader(result)

			// set content type
			ctx.Set("Content-Type", "application/octet-stream")

			// set content disposition
			ctx.Set("Content-Disposition", "attachment; filename="+asset.Name)

			return ctx.SendStream(reader, len(result))

		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
		"error": "no servercopilot release found",
	})
}

func (c *DeviceController) GetFilesystem(ctx *fiber.Ctx) error {
	versionUrl := os.Getenv("GITHUB_RELEASE_API")
	result, err := lib.GetUrl[model.GithubRelease](versionUrl, model.Unmarshal)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	for _, asset := range result.Assets {
		if asset.Name == "filesystem.bin" {
			downloadUrl := asset.BrowserDownloadURL
			result, err := lib.DownloadFile(downloadUrl, true)

			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error": err.Error(),
				})
			}

			reader := bytes.NewReader(result)

			// set content type
			ctx.Set("Content-Type", "application/octet-stream")

			// set content disposition
			ctx.Set("Content-Disposition", "attachment; filename="+asset.Name)

			return ctx.SendStream(reader, len(result))
		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
		"error": "no servercopilot release found",
	})
}

func RegisterDeviceController(c DeviceController) {
	c.Route.Get("/device/version", c.GetVersion)
	c.Route.Get("/device/firmware", c.GetFirmware)
	c.Route.Get("/device/filesystem", c.GetFilesystem)
}
