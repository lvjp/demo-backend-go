package requestid

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type contextKeyType string

var contextKey contextKeyType = "requestID"

func With(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, contextKey, requestID)
}

func MustGet(ctx context.Context) string {
	requestID, exists := ctx.Value(contextKey).(string)
	if !exists {
		panic("requestID not found in context")
	}

	return requestID
}

func Middleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestID := uuid.NewString()

		c.Set(fiber.HeaderXRequestID, requestID)
		c.SetUserContext(context.WithValue(c.UserContext(), contextKey, requestID))

		return c.Next()
	}
}
