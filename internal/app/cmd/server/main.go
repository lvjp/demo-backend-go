package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := newFiberApp()

	var serverErr error
	go func() {
		defer cancel()

		serverErr = server.Listen(":8080")
	}()

	<-ctx.Done()
	log.Println("Server shutdown sequence started")

	if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe error: %w", serverErr)
	}

	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("could not shutdown server: %w", err)
	}
	log.Println("Server terminated")

	return nil
}

func newFiberApp() *fiber.App {
	app := fiber.New()

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		scheme := "http"
		if listenData.TLS {
			scheme = "https"
		}

		log.Printf("Listening on %s://%s:%s\n", scheme, listenData.Host, listenData.Port)

		return nil
	})

	app.Hooks().OnShutdown(func() error {
		log.Println("Fiber shutdown done")
		return nil
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return app
}
