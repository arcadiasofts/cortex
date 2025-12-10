package main

import (
	"backend/config"
	"backend/internal/handlers"
	"backend/internal/services"
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Redis connect provider
func ProvideRedis(cfg *config.AppConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})

	// Connection test
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}

// ProvideAuthService Service provider
func ProvideAuthService(r *redis.Client, cfg *config.AppConfig) *services.AuthService {
	return services.NewAuthService(r, cfg.Jwt.Secret)
}

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
	authHandler *handlers.AuthHandler, // Fx가 자동으로 주입함
) {
	// Register routes
	api := app.Group("/api/v1")
	auth := api.Group("/auth")

	auth.Post("/challenge", authHandler.RequestChallenge)
	auth.Post("/login", authHandler.Login)

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

func main() {
	fx.New(
		fx.Provide(
			config.ProvideConfig,    // Configs
			ProvideRedis,            // Redis
			ProvideFiberApp,         // Fiber App
			ProvideAuthService,      // -> *services.AuthService
			handlers.NewAuthHandler, // -> *handlers.AuthHandler
		),

		fx.Invoke(RegisterAndStart),
	).Run()
}
