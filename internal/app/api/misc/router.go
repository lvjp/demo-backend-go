package misc

import (
	"github.com/gofiber/fiber/v2"
)

func Router(app fiber.Router, service Service) {
	app.Get("/version", GetVersion(service))
}
