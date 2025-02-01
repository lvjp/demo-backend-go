package handlers

import (
	"go.lvjp.me/demo-backend-go/internal/app/misc"

	"github.com/gofiber/fiber/v2"
)

func Version(service misc.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		version := service.Version()

		output := &fiber.Map{
			"revision": version.Revision,
			"time":     version.RevisionTime,
			"go":       version.GoVersion,
			"platform": version.Platform,
			"modified": version.Modified,
		}

		return c.JSON(output)
	}
}
