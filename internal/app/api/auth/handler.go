package auth

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SessionCreate(service SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input SessionCreateInput

		if err := json.Unmarshal(c.Body(), &input); err != nil {
			zerolog.Ctx(c.UserContext()).Debug().Err(err).Msg("Request body parsing error")
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if err := service.Create(input); errors.Is(err, ErrAccessDenied) {
			c.Status(fiber.StatusUnauthorized)
			return nil
		} else if err != nil {
			zerolog.Ctx(c.UserContext()).Error().Err(err).Msg("Session creation error")
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}
