package misc

import (
	"github.com/gofiber/fiber/v2"
)

func GetVersion(service Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		version := service.GetVersion()

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
