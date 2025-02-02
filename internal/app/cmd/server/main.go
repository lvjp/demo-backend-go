package server

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.lvjp.me/demo-backend-go/internal/app/api/routes"
	"go.lvjp.me/demo-backend-go/internal/app/misc"
	"go.lvjp.me/demo-backend-go/pkg/requestid"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

func Run() error {
	logger := newLogger()
	ctx := logger.WithContext(context.Background())

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger.With().Str("logger", "stdlog").Logger())

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := newFiberApp(logger)

	var serverErr error
	go func() {
		defer cancel()

		serverErr = server.Listen(":8080")
	}()

	<-ctx.Done()
	logger.Info().Msg("Server shutdown sequence started")

	if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe error: %w", serverErr)
	}

	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("could not shutdown server: %w", err)
	}
	logger.Info().Msg("Server terminated")

	return nil
}

func newFiberApp(logger *zerolog.Logger) *fiber.App {
	app := fiber.New()

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		scheme := "http"
		if listenData.TLS {
			scheme = "https"
		}

		logger.Info().
			Str("endpoint", scheme+"://"+listenData.Host+":"+listenData.Port).
			Msg("Listening")

		return nil
	})

	app.Hooks().OnShutdown(func() error {
		logger.Info().Msg("Fiber shutdown done")
		return nil
	})

	app.Use(requestid.Middleware())
	app.Use(func(c *fiber.Ctx) error {
		logger := logger.With().Str("requestID", requestid.MustGet(c.UserContext())).Logger()
		c.SetUserContext(logger.WithContext(c.UserContext()))
		return c.Next()
	})
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))

	app.Use(cors.New())

	routes.MiscRouter(app.Group("/api/v0/misc"), misc.NewService())

	return app
}

func newLogger() *zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return &logger
}
