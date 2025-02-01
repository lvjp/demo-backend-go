package routes

import (
	"go.lvjp.me/demo-backend-go/internal/app/api/handlers"
	"go.lvjp.me/demo-backend-go/internal/app/misc"

	"github.com/gofiber/fiber/v2"
)

func MiscRouter(app fiber.Router, service misc.Service) {
	app.Get("/version", handlers.Version(service))
}
