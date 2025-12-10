package main

import (
	"backend/config"
	"backend/internal/handlers"
	"backend/internal/services"
	"backend/provider/database"
	"backend/server"
	"context"
	"fmt"
	"log"

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

func main() {
	fx.New(
		fx.Provide(
			config.ProvideConfig, // Configs
			ProvideRedis,         // Redis
			database.NewDatabase,
			server.ProvideFiberApp,      // Fiber App
			services.ProvideAuthService, // -> *services.AuthService
			handlers.NewAuthHandler,     // -> *handlers.AuthHandler
		),

		fx.Invoke(server.RegisterAndStart),
	).Run()
}
