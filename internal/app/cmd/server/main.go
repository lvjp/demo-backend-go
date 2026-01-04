package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.lvjp.me/demo-backend-go/internal/app/api/auth"
	"go.lvjp.me/demo-backend-go/internal/app/api/misc"
	"go.lvjp.me/demo-backend-go/internal/app/config"
	"go.lvjp.me/demo-backend-go/internal/app/db"
	"go.lvjp.me/demo-backend-go/pkg/requestid"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

func Run() error {
	config, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load configuration: %v\n", err)
		os.Exit(1)
	}

	logger := newLogger(config.Log)
	ctx := logger.WithContext(context.Background())

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger.With().Str("logger", "stdlog").Logger())

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	conn, err := db.NewConnector(ctx, config.Database)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	server := newFiberApp(logger)

	misc.Router(server.Group("/api/v0/misc"), misc.NewService())
	auth.Router(server.Group("/api/v0/auth"), auth.NewSessionService(conn))

	var serverErr error
	go func() {
		defer cancel()

		serverErr = server.Listen(*config.Server.ListenAddress)
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

	return app
}

func newLogger(config config.Log) *zerolog.Logger {
	var writer io.Writer = os.Stderr

	if config.Format != nil {
		switch *config.Format {
		case "json":
			// default is json, do nothing
		case "console":
			writer = zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown log format %q, defaulting to json\n", *config.Format)
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}
