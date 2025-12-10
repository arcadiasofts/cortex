package server

import (
	"backend/config"
	"backend/internal/handlers"
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// ProvideFiberApp Fiber application provider
func ProvideFiberApp(
	cfg *config.AppConfig,
) *fiber.App {
	return fiber.New(fiber.Config{
		AppName: fmt.Sprintf("%s v%s", cfg.App.Name, cfg.App.Version),
	})
}

// RegisterAndStart Main invoke function.
func RegisterAndStart(
	lc fx.Lifecycle,
	app *fiber.App,
	cfg *config.AppConfig,
	authHandler *handlers.AuthHandler,
) {
	// Register routes
	api := app.Group("/api/v1")
	authHandler.Register(api)

	// Lifecycle hook
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Server starting on port %s", cfg.App.Port)
				if err := app.Listen(cfg.App.Port); err != nil {
					log.Printf("Server Error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping Cortex Server...")
			return app.Shutdown()
		},
	})
}
