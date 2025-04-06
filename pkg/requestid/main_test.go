package requestid

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	requestID := uuid.NewString()

	require.PanicsWithValue(t, "requestID not found in context", func() {
		gotRequestID := MustGet(context.Background())
		require.Empty(t, gotRequestID)
	})

	require.NotPanics(t, func() {
		ctx := With(context.Background(), requestID)
		gotRequestID := MustGet(ctx)
		require.Equal(t, requestID, gotRequestID)
	})
}

func TestMiddleware(t *testing.T) {
	var generatedID string

	app := fiber.New()
	app.Use(Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		generatedID = MustGet(c.UserContext())
		c.Status(fiber.StatusNoContent)
		return nil
	})

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	res, err := app.Test(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.NotEmpty(t, generatedID)
	require.Equal(t, fiber.StatusNoContent, res.StatusCode)
	require.Equal(t, generatedID, res.Header.Get(fiber.HeaderXRequestID))
}
