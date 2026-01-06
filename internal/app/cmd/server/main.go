package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"go.lvjp.me/demo-backend-go/internal/app/api/auth"
	"go.lvjp.me/demo-backend-go/internal/app/api/misc"
	"go.lvjp.me/demo-backend-go/internal/app/appcontext"
	"go.lvjp.me/demo-backend-go/internal/app/db"
	"go.lvjp.me/demo-backend-go/pkg/requestid"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

func Run(ctx *appcontext.AppContext) error {
	defer func() {
		ctx.Logger.Info().Msg("Server terminated")
	}()

	var cancel context.CancelFunc
	ctx.Context, cancel = signal.NotifyContext(ctx.Context, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	conn, err := db.NewConnector(ctx, ctx.Config.Database)
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			ctx.Logger.Warn().Err(err).Msg("Could not close database connection")
		}
	}()

	server := newFiberApp(&ctx.Logger)

	misc.Router(server.Group("/api/v0/misc"), misc.NewService())
	auth.Router(server.Group("/api/v0/auth"), auth.NewSessionService(conn))

	var serverErr error
	go func() {
		defer cancel()

		serverErr = server.Listen(*ctx.Config.Server.ListenAddress)
	}()

	<-ctx.Done()
	ctx.Logger.Info().Msg("Server shutdown sequence started")

	if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe error: %v", serverErr)
	}

	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("could not shutdown server: %v", err)
	}

	return nil
}

func newFiberApp(logger *zerolog.Logger) *fiber.App {
	app := fiber.New()

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		u := url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(listenData.Host, listenData.Port),
		}

		if listenData.TLS {
			u.Scheme = "https"
		}

		logger.Info().
			Stringer("endpoint", &u).
			Msg("Listening")

		return nil
	})

	app.Hooks().OnShutdown(func() error {
		logger.Info().Msg("Fiber shutdown done")
		return nil
	})

	app.Use(requestid.Middleware())

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
		Fields: slices.Concat([]string{fiberzerolog.FieldRequestID}, fiberzerolog.ConfigDefault.Fields),
	}))

	app.Use(cors.New())

	return app
}
