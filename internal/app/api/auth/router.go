package auth

import (
	"github.com/gofiber/fiber/v2"
)

func Router(app fiber.Router, service SessionService) {
	app.Post("/session", SessionCreate(service))
}
